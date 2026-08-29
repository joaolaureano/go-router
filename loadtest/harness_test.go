package loadtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/joaolaureano/go-router/loadtest/app"
)

// The knobs. Defaults are sized to run on a laptop in a few seconds; CI or a
// real soak run raises them without touching the code.
var (
	rate     = envInt("LOADTEST_RATE", 2000) // requests per second
	duration = envDur("LOADTEST_DURATION", 2*time.Second)
	workers  = envInt("LOADTEST_WORKERS", 50)
	// A p99 budget is a machine-dependent number, so it is generous by
	// default and always reported regardless. The assertions that actually
	// guard behaviour are the status and body ones below it.
	p99Budget = envDur("LOADTEST_P99", 250*time.Millisecond)
)

func envInt(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			return n
		}
	}
	return fallback
}

func envDur(key string, fallback time.Duration) time.Duration {
	if raw := os.Getenv(key); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil {
			return d
		}
	}
	return fallback
}

// server starts the routing table under test on a real socket. The load goes
// through net/http exactly as it would in production; nothing here calls
// ServeHTTP directly.
func server(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(app.New())
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)
	return srv
}

// result is one recorded hit, kept only for the fields the assertions read.
// Vegeta's own Result carries the body, which is what lets a load test check
// that the router bound the right segment and not merely that it answered.
type run struct {
	metrics vegeta.Metrics
	// byTarget maps "METHOD URL" to the distinct bodies seen for it, so a
	// scenario can assert what a route answered without holding every hit.
	byTarget map[string]map[string]int
	headers  map[string]http.Header
}

// attack fires a paced scenario at the configured rate.
func attack(t *testing.T, name string, targeter vegeta.Targeter) *run {
	t.Helper()

	return collect(t, name, targeter,
		vegeta.Rate{Freq: rate, Per: time.Second}, duration,
		vegeta.Workers(uint64(workers)),
		vegeta.KeepAlive(true),
		vegeta.Timeout(10*time.Second),
	)
}

// collect runs one attack and folds every response into a run. It builds the
// attacker itself and builds a new one per call, which is not a stylistic
// choice: vegeta's Attack closes the attacker's stop channel when it returns,
// and stopOnce makes that permanent, so a second Attack on the same Attacker
// stops after roughly one request. Everything after the attack -- the
// per-target bodies, the headers, the report -- is the same work whether the
// run was paced or a burst, which is why both go through here.
func collect(
	t *testing.T,
	name string,
	targeter vegeta.Targeter,
	pacer vegeta.Pacer,
	du time.Duration,
	opts ...func(*vegeta.Attacker),
) *run {
	t.Helper()

	attacker := vegeta.NewAttacker(opts...)
	defer attacker.Stop()

	out := &run{
		byTarget: make(map[string]map[string]int),
		headers:  make(map[string]http.Header),
	}

	for res := range attacker.Attack(targeter, pacer, du, name) {
		out.metrics.Add(res)
		key := res.Method + " " + res.URL
		bodies, ok := out.byTarget[key]
		if !ok {
			bodies = make(map[string]int)
			out.byTarget[key] = bodies
		}
		bodies[string(res.Body)]++
		if _, seen := out.headers[key]; !seen && res.Headers != nil {
			out.headers[key] = res.Headers
		}
	}
	out.metrics.Close()

	report(t, name, &out.metrics)
	return out
}

func report(t *testing.T, name string, m *vegeta.Metrics) {
	t.Helper()
	t.Logf("[%s] requests=%d rate=%.0f/s throughput=%.0f/s success=%.2f%% "+
		"mean=%s p50=%s p95=%s p99=%s max=%s",
		name, m.Requests, m.Rate, m.Throughput, m.Success*100,
		m.Latencies.Mean.Round(time.Microsecond),
		m.Latencies.P50.Round(time.Microsecond),
		m.Latencies.P95.Round(time.Microsecond),
		m.Latencies.P99.Round(time.Microsecond),
		m.Latencies.Max.Round(time.Microsecond),
	)
}

// requireNoTransportErrors fails on anything that was not an HTTP answer:
// refused connections, timeouts, resets. A 404 is a result; a reset is a bug
// or an exhausted machine, and either way the numbers below it mean nothing.
func requireNoTransportErrors(t *testing.T, out *run) {
	t.Helper()
	for _, e := range out.metrics.Errors {
		// Vegeta records non-2xx status text as an "error" too; those are the
		// scenario's business, not the transport's.
		if _, err := strconv.Atoi(e[:min(3, len(e))]); err == nil {
			continue
		}
		t.Fatalf("transport error under load: %q (errors=%v)", e, out.metrics.Errors)
	}
}

// requireStatus asserts every recorded hit carried exactly this code.
func requireStatus(t *testing.T, out *run, code uint16) {
	t.Helper()
	if len(out.metrics.StatusCodes) != 1 {
		t.Fatalf("expected every response to be %d, got %v", code, out.metrics.StatusCodes)
	}
	if got := out.metrics.StatusCodes[strconv.Itoa(int(code))]; uint64(got) != out.metrics.Requests {
		t.Fatalf("expected %d responses with status %d, got %v",
			out.metrics.Requests, code, out.metrics.StatusCodes)
	}
}

// requireBody asserts a target answered with one body and only that body. It
// is the assertion that separates "the router was fast" from "the router was
// right", and it is checked on every single hit, not on a sample.
func requireBody(t *testing.T, out *run, target, want string) {
	t.Helper()
	bodies, ok := out.byTarget[target]
	if !ok {
		t.Fatalf("no responses recorded for %s (recorded: %v)", target, keys(out.byTarget))
	}
	if len(bodies) != 1 {
		t.Fatalf("%s answered with %d distinct bodies under load: %v", target, len(bodies), bodies)
	}
	for got := range bodies {
		if got != want {
			t.Fatalf("%s answered %q, want %q", target, got, want)
		}
	}
}

func requireHeader(t *testing.T, out *run, target, key, want string) {
	t.Helper()
	h, ok := out.headers[target]
	if !ok {
		t.Fatalf("no responses recorded for %s (recorded: %v)", target, keys(out.byTarget))
	}
	if got := h.Get(key); got != want {
		t.Fatalf("%s returned %s: %q, want %q", target, key, got, want)
	}
}

func requireLatencyBudget(t *testing.T, out *run) {
	t.Helper()
	if out.metrics.Latencies.P99 > p99Budget {
		t.Errorf("p99 latency %s exceeds budget %s (raise LOADTEST_P99 if this machine is the problem)",
			out.metrics.Latencies.P99.Round(time.Microsecond), p99Budget)
	}
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// roundRobin cycles a fixed list of targets, which is what makes a scenario's
// per-target body assertions meaningful: every target is hit, repeatedly.
func roundRobin(base string, targets ...vegeta.Target) vegeta.Targeter {
	prepared := make([]vegeta.Target, len(targets))
	for i, target := range targets {
		prepared[i] = target
		prepared[i].URL = base + target.URL
		if prepared[i].Method == "" {
			prepared[i].Method = http.MethodGet
		}
	}
	return vegeta.NewStaticTargeter(prepared...)
}

func get(path string) vegeta.Target {
	return vegeta.Target{Method: http.MethodGet, URL: path}
}

func skipShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("load test: skipped under -short")
	}
	fmt.Fprintln(os.Stderr) // keep the -v log readable between scenarios
}

// requireHeaderList asserts the full ordered list of values for a repeated
// header. The middleware scenario uses it: each layer appends its own name, so
// the list is the chain the request actually walked.
func requireHeaderList(t *testing.T, out *run, target, key string, want []string) {
	t.Helper()
	h, ok := out.headers[target]
	if !ok {
		t.Fatalf("no responses recorded for %s (recorded: %v)", target, keys(out.byTarget))
	}
	got := h.Values(key)
	if len(got) != len(want) {
		t.Fatalf("%s returned %s: %v, want %v", target, key, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s returned %s: %v, want %v", target, key, got, want)
		}
	}
}

// requireStatusMix asserts the run produced only expected codes. Mixed traffic
// deliberately contains misses and rejections, so pinning one code would be
// wrong; pinning the set is not.
func requireStatusMix(t *testing.T, out *run, allowed map[int]bool) {
	t.Helper()
	for code, count := range out.metrics.StatusCodes {
		n, err := strconv.Atoi(code)
		if err != nil || !allowed[n] {
			t.Fatalf("unexpected status %q (%d responses); allowed: %v",
				code, count, keysOfInt(allowed))
		}
	}
}

func keysOfInt(m map[int]bool) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
