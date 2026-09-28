package gorouter

import (
	"net/http"

	"github.com/joaolaureano/go-router/router"
)

// Router is the surface for code that wants to accept a router rather than
// construct one. It is not meant to mirror every method *router.Router offers:
// NewRouter returns the concrete type, so new methods can be added there
// without widening this. The assertion below keeps the two from drifting apart
// in the direction that matters.
type Router interface {
	http.Handler

	// Register is a method for adding a new route to the router.
	// It takes an HTTP method, a path, and a handler function as parameters.
	// The router will use these parameters to associate incoming requests with the specified handler.
	Register(httpMethod router.Method, path string, method http.HandlerFunc)

	// Use is a method for adding middleware to the router.
	// Middleware functions can process or modify requests before reaching the route handler.
	// This method enhances the router's functionality by allowing the insertion of additional processing steps.
	Use(middleware func(http.Handler) http.Handler)

	// NotFound sets the handler for routes that are not found.
	// It takes a http.HandlerFunc as a parameter and assigns it as the handler for 404 routes.
	NotFound(notFoundFn http.HandlerFunc)

	// MethodNotAllowed sets the handler for requests whose path exists but is
	// not registered under the requested method. The Allow header is already
	// set by the time it runs.
	MethodNotAllowed(methodNotAllowedFn http.HandlerFunc)

	// Mount grafts another router's routes in under a prefix, handlers and
	// middleware intact.
	Mount(prefix string, other *router.Router)

	// Group creates a subgroup of routes with a common prefix.
	// It takes a prefix string and a function that operates on a router.Router as parameters.
	// This method allows organizing routes under a shared path prefix.
	Group(prefix string, fn func(r *router.Router)) *router.Router

	// With adds one or multiple middleware functions to a specific set of routes.
	// It takes one or more middleware functions as parameters and returns a pointer to the router.Router.
	// This method provides a way to apply middleware to a subset of routes.
	With(middleware ...func(http.Handler) http.Handler) *router.Router
}

// The routing vocabulary, re-exported so that a caller importing only this
// package can name a method.
type Method = router.Method

const (
	GET     = router.GET
	HEAD    = router.HEAD
	POST    = router.POST
	PUT     = router.PUT
	PATCH   = router.PATCH
	DELETE  = router.DELETE
	OPTIONS = router.OPTIONS
)

var _ Router = (*router.Router)(nil)

// NewRouter returns the concrete router. Returning the struct rather than the
// interface is what lets Router stay a small accept-surface.
func NewRouter() *router.Router {
	return router.NewRouter()
}

// NewPrefixRouter returns a router that registers every route under prefix.
func NewPrefixRouter(prefix string) *router.Router {
	return router.NewPrefixRouter(prefix)
}
