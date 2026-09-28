package loadtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/joaolaureano/go-router/loadtest/app"
)

// The scenarios in scenarios_test.go are paced: a fixed rate against a small,
// hand-written set of routes. These are the other half -- a full route table
// with nested subroutes, hit all at once with no pacing at all.
//
// Unpaced is a different test, not just a faster one. At a fixed rate the
// harness decides when the next request goes out, so the router is never the
// thing setting the pace; with the pacer removed, every worker sends its next
// request the instant the previous answer lands, and the tree is what the
// throughput is measuring.

var (
	burstWorkers  = envInt("LOADTEST_BURST_WORKERS", 200)
	burstDuration = envDur("LOADTEST_BURST_DURATION", 3*time.Second)
)

// burstOpts caps the attacker at a fixed pool of connections. Without a
// maximum, an unpaced attack spawns a worker per outstanding request and ends
// up measuring how many goroutines the generator can hold rather than how fast
// the router answers.
func burstOpts() []func(*vegeta.Attacker) {
	return []func(*vegeta.Attacker){
		vegeta.Workers(uint64(burstWorkers)),
		vegeta.MaxWorkers(uint64(burstWorkers)),
		vegeta.KeepAlive(true),
		vegeta.Timeout(30 * time.Second),
	}
}

// burst attacks with the zero-value Rate, which vegeta reads as "no pacing":
// requests go out as fast as burstWorkers connections can carry them.
func burst(t *testing.T, name string, targeter vegeta.Targeter) *run {
	t.Helper()
	return collect(t, name, targeter, vegeta.Rate{}, burstDuration, burstOpts()...)
}

// TestBurst_EveryRouteInTheTable is the usage test: build the whole table --
// three API versions as nested groups, eight resources each, sub-collections,
// static siblings, mounted subtrees -- then hammer every single route in it
// with no pacing, and check that each one answered its own body.
//
// The pass condition is not throughput. It is that after a few hundred
// thousand unpaced requests spread over the entire tree, every route still
// answered exactly the pattern it was registered under and exactly the
// segments it should have bound. That is what catches a request landing on a
// neighbouring node, which is the failure a wide tree under burst can produce
// and a single-route benchmark never will.
func TestBurst_EveryRouteInTheTable(t *testing.T) {
	skipShort(t)

	r, routes := app.NewCatalog()
	srv := httptest.NewServer(r)
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)

	targets := make([]vegeta.Target, len(routes))
	for i, route := range routes {
		targets[i] = vegeta.Target{
			Method: string(route.Method),
			URL:    srv.URL + route.Request,
		}
	}
	t.Logf("route table: %d routes across %d nested groups and %d mounted subtrees",
		len(routes), 3, 3)

	out := burst(t, "burst:every-route", vegeta.NewStaticTargeter(targets...))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)

	// Every route in the table, checked on every hit it received. Nothing is
	// sampled and nothing is skipped: a route the burst never reached is a
	// failure too, because the table and the targets came from the same list.
	unreached := 0
	for _, route := range routes {
		target := string(route.Method) + " " + srv.URL + route.Request
		if _, ok := out.byTarget[target]; !ok {
			unreached++
			t.Errorf("route %s %s was never reached by the burst", route.Method, route.Pattern)
			continue
		}
		requireBody(t, out, target, route.Want)
	}
	if unreached == 0 {
		t.Logf("all %d routes reached and verified on every hit (%d requests)",
			len(routes), out.metrics.Requests)
	}
}

// A burst with no warm-up: the tree is built and attacked immediately, so the
// first thousands of requests arrive with nothing cached, no connection pool
// established and the tree freshly swapped in. Routing must be correct from
// the first request, not from the hundredth.
func TestBurst_ColdTree(t *testing.T) {
	skipShort(t)

	r, routes := app.NewCatalog()
	srv := httptest.NewServer(r)
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)

	// Deliberately no warm-up request here.
	targets := make([]vegeta.Target, 0, 64)
	for _, route := range routes[:min(64, len(routes))] {
		targets = append(targets, vegeta.Target{
			Method: string(route.Method),
			URL:    srv.URL + route.Request,
		})
	}

	out := collect(t, "burst:cold-tree", vegeta.NewStaticTargeter(targets...),
		vegeta.Rate{}, time.Second, burstOpts()...)

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	for _, route := range routes[:min(64, len(routes))] {
		requireBody(t, out, string(route.Method)+" "+srv.URL+route.Request, route.Want)
	}
}

// A spike: idle, then a jump straight to full rate, which is the shape of real
// traffic far more often than a constant rate is. The router has no queue and
// no admission control, so what this checks is that the spike is absorbed by
// the connection layer and not turned into wrong answers.
func TestBurst_Spike(t *testing.T) {
	skipShort(t)

	r, routes := app.NewCatalog()
	srv := httptest.NewServer(r)
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)

	// A handful of routes from different corners of the tree, so the spike is
	// spread across branches rather than concentrated on one node.
	picked := pick(routes, 12)
	targets := make([]vegeta.Target, len(picked))
	for i, route := range picked {
		targets[i] = vegeta.Target{Method: string(route.Method), URL: srv.URL + route.Request}
	}
	targeter := vegeta.NewStaticTargeter(targets...)

	// The trickle phase and the spike phase each get their own attacker,
	// because vegeta's is single-use. That means the spike opens fresh
	// connections rather than inheriting warm ones -- which makes this the
	// harsher of the two readings, and the one worth having.
	idle := collect(t, "spike:idle", targeter,
		vegeta.Rate{Freq: 50, Per: time.Second}, time.Second, burstOpts()...)
	requireNoTransportErrors(t, idle)
	requireStatus(t, idle, http.StatusOK)

	spike := collect(t, "spike:full", targeter, vegeta.Rate{}, 2*time.Second, burstOpts()...)
	requireNoTransportErrors(t, spike)
	requireStatus(t, spike, http.StatusOK)
	for _, route := range picked {
		requireBody(t, spike, string(route.Method)+" "+srv.URL+route.Request, route.Want)
	}

	factor := spike.metrics.Throughput / idle.metrics.Throughput
	t.Logf("spike absorbed: %.0f/s -> %.0f/s (%.0fx), p99 %s -> %s",
		idle.metrics.Throughput, spike.metrics.Throughput, factor,
		idle.metrics.Latencies.P99.Round(time.Microsecond),
		spike.metrics.Latencies.P99.Round(time.Microsecond))
}

// A ramp from a trickle to the machine's ceiling, reported in slices. It
// asserts correctness throughout; the numbers are there to show where latency
// starts to bend, which is the point at which the load generator, the socket
// layer or the router stopped keeping up -- and which of the three it was is
// exactly what the shape of the curve tells you.
func TestBurst_Ramp(t *testing.T) {
	skipShort(t)

	r, routes := app.NewCatalog()
	srv := httptest.NewServer(r)
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)

	picked := pick(routes, 12)
	targets := make([]vegeta.Target, len(picked))
	for i, route := range picked {
		targets[i] = vegeta.Target{Method: string(route.Method), URL: srv.URL + route.Request}
	}
	targeter := vegeta.NewStaticTargeter(targets...)

	type step struct {
		rate       int
		throughput float64
		p50, p99   time.Duration
	}
	var steps []step

	for _, r := range []int{1000, 5000, 10000, 20000, 40000} {
		out := collect(t, fmt.Sprintf("ramp:%d", r), targeter,
			vegeta.Rate{Freq: r, Per: time.Second}, time.Second, burstOpts()...)
		requireNoTransportErrors(t, out)
		requireStatus(t, out, http.StatusOK)
		steps = append(steps, step{r, out.metrics.Throughput,
			out.metrics.Latencies.P50, out.metrics.Latencies.P99})
	}

	t.Log("ramp:")
	t.Logf("  %10s %14s %10s %10s", "offered", "throughput", "p50", "p99")
	for _, s := range steps {
		t.Logf("  %8d/s %12.0f/s %10s %10s", s.rate, s.throughput,
			s.p50.Round(time.Microsecond), s.p99.Round(time.Microsecond))
	}
}

// pick spreads a selection across the table instead of taking a contiguous
// slice, which would land entirely inside one resource's subtree.
func pick(routes []app.Route, n int) []app.Route {
	if n >= len(routes) {
		return routes
	}
	stride := len(routes) / n
	picked := make([]app.Route, 0, n)
	for i := 0; len(picked) < n; i += stride {
		picked = append(picked, routes[i])
	}
	sort.Slice(picked, func(i, j int) bool { return picked[i].Pattern < picked[j].Pattern })
	return picked
}
