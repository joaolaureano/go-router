package loadtest

import (
	"net/http"
	"testing"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Each scenario isolates one decision the tree makes, so a regression names
// itself instead of hiding in a mixed average. The mixed scenario at the end
// is the one that says whether they hold together.

// A depth-1 static hit: the cheapest path through the tree, and the number
// every other scenario should be read against.
func TestLoad_StaticRoute(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "static", roundRobin(srv.URL, get("/ping")))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/ping", "pong")
	requireLatencyBudget(t, out)
}

// One capture. The body carries the bound value, so a router that matched the
// right node but bound the wrong segment fails here rather than passing.
func TestLoad_SingleParameter(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "param-1", roundRobin(srv.URL,
		get("/users/42"),
		get("/users/7"),
		get("/users/a-long-identifier-with-dashes"),
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/users/42", "id=42")
	requireBody(t, out, "GET "+srv.URL+"/users/7", "id=7")
	requireBody(t, out, "GET "+srv.URL+"/users/a-long-identifier-with-dashes",
		"id=a-long-identifier-with-dashes")
	requireLatencyBudget(t, out)
}

// Two captures with static segments between them: the walk has to leave a
// parameter node, match static children, and capture again.
func TestLoad_NestedParameters(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "param-2", roundRobin(srv.URL,
		get("/users/42/posts/1001"),
		get("/api/v1/items/abc/reviews/xyz"),
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/users/42/posts/1001", "id=42;postID=1001")
	requireBody(t, out, "GET "+srv.URL+"/api/v1/items/abc/reviews/xyz", "itemID=abc;reviewID=xyz")
	requireLatencyBudget(t, out)
}

// A static route and a parameter route that both match the same request.
// Precedence is a correctness property, and it has to hold at rate too.
func TestLoad_StaticBeatsParameterUnderLoad(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "static-vs-param", roundRobin(srv.URL,
		get("/users/me"),
		get("/users/notme"),
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/users/me", "me")
	requireBody(t, out, "GET "+srv.URL+"/users/notme", "id=notme")
	requireLatencyBudget(t, out)
}

// Middleware is per-request work, so the chain is worth measuring on its own:
// a group two levels deep, and an explicit five-layer With stack.
func TestLoad_MiddlewareChain(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "middleware", roundRobin(srv.URL,
		get("/api/health"),
		get("/api/v1/items/99"),
		get("/deep/99"),
		get("/admin/stats"),
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/api/health", "ok")
	requireBody(t, out, "GET "+srv.URL+"/api/v1/items/99", "itemID=99")
	requireBody(t, out, "GET "+srv.URL+"/deep/99", "id=99")
	requireBody(t, out, "GET "+srv.URL+"/admin/stats", "stats")

	// Every layer the request passed through left its mark, in order.
	requireHeaderList(t, out, "GET "+srv.URL+"/api/v1/items/99", "X-Chain", []string{"api", "v1"})
	requireHeaderList(t, out, "GET "+srv.URL+"/deep/99", "X-Chain",
		[]string{"d1", "d2", "d3", "d4", "d5"})
	requireHeaderList(t, out, "GET "+srv.URL+"/admin/stats", "X-Chain", []string{"admin"})
	requireLatencyBudget(t, out)
}

// A mounted subtree keeps its own handlers and middleware after the graft.
func TestLoad_MountedSubtree(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "mount", roundRobin(srv.URL,
		get("/admin/users/17"),
		vegeta.Target{Method: http.MethodDelete, URL: "/admin/users/17"},
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/admin/users/17", "id=17")
	requireBody(t, out, "DELETE "+srv.URL+"/admin/users/17", "id=17")
	requireLatencyBudget(t, out)
}

// The miss path costs whatever a full failed walk costs, and it is the path
// an internet-facing router walks most under scanning traffic.
func TestLoad_NotFound(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "404", roundRobin(srv.URL,
		get("/nope"),
		get("/users/42/posts"),                 // a prefix of a real route
		get("/users/42/posts/1001/comments/9"), // one segment past a real route
		get("/api/v2/items/1"),                 // a group that does not exist
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusNotFound)
	requireBody(t, out, "GET "+srv.URL+"/nope", "not found")
	requireLatencyBudget(t, out)
}

// A known path under an unregistered method: 405 plus a rendered Allow header,
// which is string building on the hot path and therefore worth its own run.
func TestLoad_MethodNotAllowed(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "405", roundRobin(srv.URL,
		vegeta.Target{Method: http.MethodDelete, URL: "/users/42/posts/1001"},
		vegeta.Target{Method: http.MethodPatch, URL: "/users/42/avatar"},
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusMethodNotAllowed)
	// GET and POST are registered; HEAD and OPTIONS the router answers itself.
	requireHeader(t, out, "DELETE "+srv.URL+"/users/42/posts/1001", "Allow",
		"GET, HEAD, OPTIONS, POST")
	requireHeader(t, out, "PATCH "+srv.URL+"/users/42/avatar", "Allow",
		"OPTIONS, PUT")
	requireLatencyBudget(t, out)
}

// HEAD with no route of its own falls back to GET, which costs a second
// lookup on every such request.
func TestLoad_HeadFallsBackToGet(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "head-fallback", roundRobin(srv.URL,
		vegeta.Target{Method: http.MethodHead, URL: "/ping"},
		vegeta.Target{Method: http.MethodHead, URL: "/users/42"},
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	// net/http suppresses the body for HEAD; the point is that it routed.
	requireBody(t, out, "HEAD "+srv.URL+"/ping", "")
	requireLatencyBudget(t, out)
}

// OPTIONS on a known path with nothing registered for it: 204 and an Allow
// header, answered by the router rather than by a handler.
func TestLoad_OptionsAnsweredAutomatically(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "options", roundRobin(srv.URL,
		vegeta.Target{Method: http.MethodOptions, URL: "/ping"},
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusNoContent)
	requireHeader(t, out, "OPTIONS "+srv.URL+"/ping", "Allow", "GET, HEAD, OPTIONS")
	requireLatencyBudget(t, out)
}

// Percent-escaping is decoded per segment, so an escaped slash stays inside
// one variable instead of becoming a separator -- and only escaped requests
// pay for the decode.
func TestLoad_PercentEncodedSegments(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "escaped", roundRobin(srv.URL,
		get("/users/a%2Fb"),
		get("/users/caf%C3%A9"),
	))

	requireNoTransportErrors(t, out)
	requireStatus(t, out, http.StatusOK)
	requireBody(t, out, "GET "+srv.URL+"/users/a%2Fb", "id=a/b")
	requireBody(t, out, "GET "+srv.URL+"/users/caf%C3%A9", "id=café")
	requireLatencyBudget(t, out)
}

// The shape real traffic has: every branch above, interleaved, on one tree.
func TestLoad_MixedTraffic(t *testing.T) {
	skipShort(t)
	srv := server(t)

	out := attack(t, "mixed", roundRobin(srv.URL,
		get("/ping"),
		get("/users/42"),
		get("/users/me"),
		get("/users/42/posts/1001"),
		get("/api/health"),
		get("/api/v1/items/99/reviews/7"),
		get("/admin/users/17"),
		get("/deep/5"),
		get("/users/a%2Fb"),
		get("/nope"),
		vegeta.Target{Method: http.MethodDelete, URL: "/users/42/posts/1001"},
		vegeta.Target{Method: http.MethodOptions, URL: "/ping"},
		vegeta.Target{Method: http.MethodHead, URL: "/ping"},
	))

	requireNoTransportErrors(t, out)
	requireStatusMix(t, out, map[int]bool{
		http.StatusOK:               true,
		http.StatusNoContent:        true,
		http.StatusNotFound:         true,
		http.StatusMethodNotAllowed: true,
	})
	requireBody(t, out, "GET "+srv.URL+"/ping", "pong")
	requireBody(t, out, "GET "+srv.URL+"/users/42/posts/1001", "id=42;postID=1001")
	requireBody(t, out, "GET "+srv.URL+"/nope", "not found")
	requireLatencyBudget(t, out)
}
