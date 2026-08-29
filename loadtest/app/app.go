// Package app builds the router the load tests attack.
//
// It exists so the in-process vegeta suite and the standalone server the
// vegeta CLI shoots at are provably the same routing table: a number measured
// against one is a number about the other.
package app

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/joaolaureano/go-router/router"
	"github.com/joaolaureano/go-router/router/context"
)

// Requests counts every request that reached a handler. The load tests use it
// to tell "the router answered" apart from "something answered 200".
var Requests atomic.Int64

func text(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(body))
}

// echo answers with the route's captured variables so that a load test can
// catch a router that routes fast but binds the wrong segment.
func echo(names ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Requests.Add(1)
		body := ""
		for i, name := range names {
			if i > 0 {
				body += ";"
			}
			body += name + "=" + context.Param(r, name)
		}
		text(w, body)
	}
}

func static(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Requests.Add(1)
		text(w, body)
	}
}

// tag is the middleware the load tests use to prove the chain still runs at
// rate: every layer it passes through leaves a header behind.
func tag(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Chain", name)
			next.ServeHTTP(w, r)
		})
	}
}

// New builds the routing table under load. Every route here corresponds to a
// distinct decision inside the tree -- static hit, parameter capture, deeper
// static under a parameter, catch-all, mounted subtree -- so that a regression
// in any one of them shows up as its own scenario rather than as a blurred
// average.
func New() *router.Router {
	r := router.NewRouter()

	// Depth-1 static: the cheapest lookup the tree can do, and the baseline
	// every other scenario is read against.
	r.Get("/ping", static("pong"))

	// Single capture.
	r.Get("/users/{id}", echo("id"))

	// Capture, then static, then capture: forces the walk past a parameter
	// node and back into static children.
	r.Get("/users/{id}/posts/{postID}", echo("id", "postID"))

	// A static sibling of the parameter above. Static must win at rate, not
	// just in a unit test.
	r.Get("/users/me", static("me"))

	// Registered under POST only, on a path whose GET sibling exists: the
	// 405 scenario, and the Allow header it has to render.
	r.Post("/users/{id}/posts/{postID}", echo("id", "postID"))

	// A path that exists under exactly one method, used for the 405 scenario
	// where nothing else shares the node.
	r.Put("/users/{id}/avatar", echo("id"))

	// Catch-all: the tail lands in one variable, slashes included.
	r.Get("/files/*", echo(router.WildcardParam))

	// Group with its own middleware stack, nested one level deeper. Three
	// layers of chain on every request that lands here.
	r.Group("/api", func(api *router.Router) {
		api.Use(tag("api"))
		api.Get("/health", static("ok"))

		api.Group("/v1", func(v1 *router.Router) {
			v1.Use(tag("v1"))
			v1.Get("/items/{itemID}", echo("itemID"))
			v1.Get("/items/{itemID}/reviews/{reviewID}", echo("itemID", "reviewID"))
		})
	})

	// A mounted subtree, middleware and all, grafted under a prefix.
	r.Mount("/admin", adminRouter())

	// Deliberate deep chain: how a request pays for middleware it has to walk
	// through before reaching a handler.
	deep := r.With(tag("d1"), tag("d2"), tag("d3"), tag("d4"), tag("d5"))
	deep.Get("/deep/{id}", echo("id"))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		Requests.Add(1)
		w.WriteHeader(http.StatusNotFound)
		text(w, "not found")
	})

	return r
}

func adminRouter() *router.Router {
	admin := router.NewRouter()
	admin.Use(tag("admin"))
	admin.Get("/stats", static("stats"))
	admin.Get("/users/{id}", echo("id"))
	admin.Delete("/users/{id}", echo("id"))
	return admin
}

// Describe lists the routes for the standalone server's banner.
func Describe() string {
	return fmt.Sprint(
		"GET    /ping\n",
		"GET    /users/me\n",
		"GET    /users/{id}\n",
		"GET    /users/{id}/posts/{postID}\n",
		"POST   /users/{id}/posts/{postID}\n",
		"PUT    /users/{id}/avatar\n",
		"GET    /files/*\n",
		"GET    /api/health\n",
		"GET    /api/v1/items/{itemID}\n",
		"GET    /api/v1/items/{itemID}/reviews/{reviewID}\n",
		"GET    /admin/stats\n",
		"GET    /admin/users/{id}\n",
		"DELETE /admin/users/{id}\n",
		"GET    /deep/{id}\n",
	)
}
