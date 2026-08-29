package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/joaolaureano/go-router/router"
	"github.com/joaolaureano/go-router/router/context"
)

// A burst against a route table is only worth running if it covers the whole
// table, and it can only assert on the answers if it knows what each route
// should say. So the table and the expectations come from one place: Catalog
// builds the route list, NewCatalog registers exactly that list, and the burst
// test turns the same list into targets. A route that exists but is never hit
// is impossible by construction.

// Route is one entry: how it is registered, one concrete request that must
// reach it, and the exact answer that request must get.
type Route struct {
	Method  router.Method
	Pattern string // as registered, with {variables}
	Request string // a concrete path that must land on Pattern
	Want    string // the exact body the handler answers
	Status  int
}

// Target renders the "METHOD path" key the load harness records results under.
func (route Route) Target() string {
	return string(route.Method) + " " + route.Request
}

// identity is what every catalog handler answers: the pattern it was
// registered under, followed by the variables it bound. It is the assertion
// that a burst needs -- under a few hundred routes and no pacing, the failure
// worth catching is not a slow response but a request landing on the wrong
// node, or on the right node with the wrong segment bound. Either shows up
// here as a body that does not match.
func identity(pattern string) http.HandlerFunc {
	names := variablesOf(pattern)
	return func(w http.ResponseWriter, r *http.Request) {
		Requests.Add(1)
		var body strings.Builder
		body.WriteString(pattern)
		body.WriteByte('|')
		for i, name := range names {
			if i > 0 {
				body.WriteByte(';')
			}
			body.WriteString(name)
			body.WriteByte('=')
			body.WriteString(context.Param(r, name))
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(body.String()))
	}
}

// variablesOf lists a pattern's capture names in order, the catch-all
// included under the name it binds to.
func variablesOf(pattern string) []string {
	var names []string
	for _, segment := range strings.Split(strings.TrimPrefix(pattern, "/"), "/") {
		switch {
		case segment == router.WildcardParam:
			names = append(names, router.WildcardParam)
		case strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}"):
			names = append(names, strings.Trim(segment, "{}"))
		}
	}
	return names
}

// wantFor computes the body a request to pattern must produce, by binding the
// concrete segments of request to the pattern's variables. It is deliberately
// a second, independent walk of the two strings: if it agreed with the router
// by sharing its code, it would agree with the router's bugs too.
func wantFor(pattern, request string) string {
	patternSegments := strings.Split(strings.TrimPrefix(pattern, "/"), "/")
	requestSegments := strings.Split(strings.TrimPrefix(request, "/"), "/")

	var bound []string
	for i, segment := range patternSegments {
		switch {
		case segment == router.WildcardParam:
			// The catch-all takes the whole remaining tail, slashes included.
			bound = append(bound, router.WildcardParam+"="+strings.Join(requestSegments[i:], "/"))
		case strings.HasPrefix(segment, "{"):
			bound = append(bound, strings.Trim(segment, "{}")+"="+requestSegments[i])
		}
	}
	return pattern + "|" + strings.Join(bound, ";")
}

// The shape of the generated table: three API versions, each a nested group,
// each carrying the same resources, each resource with an id, a static sibling
// of that id, two sub-collections under it and a catch-all. It is the shape a
// real service's routing table has -- wide at the resource level, three to
// five deep, static and parameter siblings at the same node -- rather than a
// microbenchmark's single branch.
var (
	catalogVersions  = []string{"v1", "v2", "v3"}
	catalogResources = []string{"users", "orders", "products", "invoices", "teams", "projects", "events", "webhooks"}
	catalogSubs      = []string{"comments", "tags"}
)

// Catalog returns every route in the generated table, with a concrete request
// and expected answer for each.
func Catalog() []Route {
	var routes []Route

	add := func(method router.Method, pattern, request string) {
		routes = append(routes, Route{
			Method:  method,
			Pattern: pattern,
			Request: request,
			Want:    wantFor(pattern, request),
			Status:  http.StatusOK,
		})
	}

	for _, version := range catalogVersions {
		for _, resource := range catalogResources {
			base := fmt.Sprintf("/api/%s/%s", version, resource)
			id := resource + "-42"

			add(router.GET, base, base)
			add(router.POST, base, base)

			// A static child sitting next to the parameter child: precedence
			// has to hold at every one of these nodes, not just at one.
			add(router.GET, base+"/count", base+"/count")

			add(router.GET, base+"/{id}", base+"/"+id)
			add(router.PUT, base+"/{id}", base+"/"+id)
			add(router.DELETE, base+"/{id}", base+"/"+id)

			for _, sub := range catalogSubs {
				add(router.GET, base+"/{id}/"+sub, base+"/"+id+"/"+sub)
				add(router.GET, base+"/{id}/"+sub+"/{subID}", base+"/"+id+"/"+sub+"/"+sub+"-7")
			}

			// A catch-all deep inside the tree, not at the root: the walk has
			// to get past four levels before it can fall back to this.
			add(router.GET, base+"/{id}/files/*", base+"/"+id+"/files/deep/nested/report.pdf")
		}
	}

	// Mounted subtrees, one per version, grafted after the fact.
	for _, version := range catalogVersions {
		base := "/mounted/" + version
		add(router.GET, base+"/status", base+"/status")
		add(router.GET, base+"/{id}", base+"/m-9")
		add(router.GET, base+"/{id}/detail", base+"/m-9/detail")
	}

	return routes
}

// NewCatalog builds a router registering exactly the routes Catalog returns,
// and returns both so a test can assert over the same list it registered.
func NewCatalog() (*router.Router, []Route) {
	routes := Catalog()
	r := router.NewRouter()

	// The versions are nested groups rather than flat registrations: subroutes
	// under a shared prefix are what the table is meant to exercise, and a
	// group's prefix composition is part of what the burst is testing.
	r.Group("/api", func(api *router.Router) {
		for _, version := range catalogVersions {
			api.Group("/"+version, func(v *router.Router) {
				prefix := "/api/" + version
				for _, route := range routes {
					if !strings.HasPrefix(route.Pattern, prefix+"/") && route.Pattern != prefix {
						continue
					}
					v.Register(route.Method, strings.TrimPrefix(route.Pattern, prefix), identity(route.Pattern))
				}
			})
		}
	})

	// The mounted trees are built whole and grafted, which is a single large
	// swap rather than one registration per route.
	for _, version := range catalogVersions {
		prefix := "/mounted/" + version
		sub := router.NewRouter()
		for _, route := range routes {
			if !strings.HasPrefix(route.Pattern, prefix+"/") {
				continue
			}
			sub.Register(route.Method, strings.TrimPrefix(route.Pattern, prefix), identity(route.Pattern))
		}
		r.Mount(prefix, sub)
	}

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		Requests.Add(1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	return r, routes
}
