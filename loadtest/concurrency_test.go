package loadtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/joaolaureano/go-router/loadtest/app"
	"github.com/joaolaureano/go-router/router"
)

// The router swaps its tree atomically so that registration does not block
// serving. A unit test can show that is race-free; only load can show it stays
// correct while a real socket is saturated, which is what these two do.

// Routes registered mid-flight must appear without disturbing the ones already
// being served. The attack runs against a route that exists throughout, while
// a writer keeps grafting new ones underneath it.
func TestLoad_RegistrationWhileServing(t *testing.T) {
	skipShort(t)

	r := app.New()
	srv := httptest.NewServer(r)
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)

	stop := make(chan struct{})
	var writer sync.WaitGroup
	writer.Add(1)
	registered := 0
	go func() {
		defer writer.Done()
		for i := 0; ; i++ {
			select {
			case <-stop:
				registered = i
				return
			default:
			}
			path := fmt.Sprintf("/generated/%d/{id}", i)
			r.Get(path, func(w http.ResponseWriter, req *http.Request) {
				w.Write([]byte("generated"))
			})
			// Fast enough to overlap the attack many times over, slow enough
			// that the test measures the router and not the scheduler.
			time.Sleep(200 * time.Microsecond)
		}
	}()

	out := attack(t, "register-while-serving", roundRobin(srv.URL,
		get("/ping"),
		get("/users/42"),
		get("/users/42/posts/1001"),
		get("/api/v1/items/99"),
	))

	close(stop)
	writer.Wait()
	t.Logf("registered %d routes during the attack", registered)

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	// The pre-existing routes answered exactly as before, on every hit.
	requireBody(t, out, "GET "+srv.URL+"/ping", "pong")
	requireBody(t, out, "GET "+srv.URL+"/users/42", "id=42")
	requireBody(t, out, "GET "+srv.URL+"/users/42/posts/1001", "id=42;postID=1001")
	requireBody(t, out, "GET "+srv.URL+"/api/v1/items/99", "itemID=99")
	requireLatencyBudget(t, out)

	if registered == 0 {
		t.Fatal("no routes were registered during the attack; the scenario proved nothing")
	}

	// And the last route the writer added is reachable afterwards, so the
	// swaps were not merely harmless -- they landed.
	resp, err := http.Get(fmt.Sprintf("%s/generated/%d/7", srv.URL, registered-1))
	if err != nil {
		t.Fatalf("route registered under load is unreachable: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("route registered under load answered %d, want 200", resp.StatusCode)
	}
}

// Mount grafts a whole subtree at once, which is a larger swap than a single
// registration. Traffic in flight across the graft must never see a half-built
// tree: no 404 for a route that existed before it, and none for one the mount
// brought in once it has landed.
func TestLoad_MountWhileServing(t *testing.T) {
	skipShort(t)

	r := app.New()
	srv := httptest.NewServer(r)
	srv.Config.SetKeepAlivesEnabled(true)
	t.Cleanup(srv.Close)

	stop := make(chan struct{})
	var writer sync.WaitGroup
	writer.Add(1)
	mounted := 0
	go func() {
		defer writer.Done()
		for i := 0; ; i++ {
			select {
			case <-stop:
				mounted = i
				return
			default:
			}
			sub := router.NewRouter()
			sub.Get("/thing", func(w http.ResponseWriter, req *http.Request) {
				w.Write([]byte("mounted"))
			})
			sub.Get("/thing/{id}", func(w http.ResponseWriter, req *http.Request) {
				w.Write([]byte("mounted"))
			})
			r.Mount(fmt.Sprintf("/mounted/%d", i), sub)
			time.Sleep(500 * time.Microsecond)
		}
	}()

	out := attack(t, "mount-while-serving", roundRobin(srv.URL,
		get("/ping"),
		get("/admin/users/17"),
		get("/api/v1/items/99"),
	))

	close(stop)
	writer.Wait()
	t.Logf("mounted %d subtrees during the attack", mounted)

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/ping", "pong")
	requireBody(t, out, "GET "+srv.URL+"/admin/users/17", "id=17")
	requireLatencyBudget(t, out)

	if mounted == 0 {
		t.Fatal("nothing was mounted during the attack; the scenario proved nothing")
	}
}

// A comparative run: the same rate against each route shape, reported side by
// side. The absolute numbers are machine-dependent and assert nothing; what is
// worth reading is the distance between the rows, which is where the cost of a
// parameter, a decode or a middleware layer becomes visible.
func TestLoad_ShapeComparison(t *testing.T) {
	skipShort(t)
	srv := server(t)

	shapes := []struct {
		name   string
		target vegeta.Target
	}{
		{"static depth 1", get("/ping")},
		{"static depth 2", get("/users/me")},
		{"1 parameter", get("/users/42")},
		{"2 parameters", get("/users/42/posts/1001")},
		{"escaped segment", get("/users/a%2Fb")},
		{"2 middleware", get("/api/v1/items/99")},
		{"5 middleware", get("/deep/99")},
		{"mounted subtree", get("/admin/users/17")},
		{"404 (full miss)", get("/nope")},
		{"404 (deep miss)", get("/users/42/posts/1001/comments/9")},
		{"405 + Allow", vegeta.Target{Method: http.MethodDelete, URL: "/users/42/posts/1001"}},
		{"HEAD -> GET", vegeta.Target{Method: http.MethodHead, URL: "/ping"}},
		{"OPTIONS auto", vegeta.Target{Method: http.MethodOptions, URL: "/ping"}},
	}

	type row struct {
		name       string
		throughput float64
		mean       time.Duration
		p99        time.Duration
	}
	rows := make([]row, 0, len(shapes))

	for _, shape := range shapes {
		out := attack(t, shape.name, roundRobin(srv.URL, shape.target))
		requireNoTransportErrors(t, out)
		rows = append(rows, row{
			name:       shape.name,
			throughput: out.metrics.Throughput,
			mean:       out.metrics.Latencies.Mean,
			p99:        out.metrics.Latencies.P99,
		})
	}

	t.Log("route shape comparison at " + fmt.Sprint(rate) + " req/s:")
	t.Logf("  %-26s %12s %12s %12s", "shape", "throughput", "mean", "p99")
	for _, r := range rows {
		t.Logf("  %-26s %10.0f/s %12s %12s",
			r.name, r.throughput,
			r.mean.Round(time.Microsecond), r.p99.Round(time.Microsecond))
	}
}
