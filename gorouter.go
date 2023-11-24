package main

import (
	"net/http"

	_const "web/const"
	"web/router"
)

// Router is an interface that extends the http.Handler interface,
// requiring any implementing type to fulfill the http.Handler contract.
// It declares two methods - Register and Use - for registering routes and adding middleware.
type Router interface {
	http.Handler

	// Register is a method for adding a new route to the router.
	// It takes an HTTP method, a path, and a handler function as parameters.
	// The router will use these parameters to associate incoming requests with the specified handler.
	Register(httpMethod _const.HTTPMethods, path string, method http.HandlerFunc)

	// Use is a method for adding middleware to the router.
	// Middleware functions can process or modify requests before reaching the route handler.
	// This method enhances the router's functionality by allowing the insertion of additional processing steps.
	Use(middleware func(http.Handler) http.Handler)

	// NotFound sets the handler for routes that are not found.
	// It takes a http.HandlerFunc as a parameter and assigns it as the handler for 404 routes.
	NotFound(notFoundFn http.HandlerFunc)

	// Group creates a subgroup of routes with a common prefix.
	// It takes a prefix string and a function that operates on a router.Router as parameters.
	// This method allows organizing routes under a shared path prefix.
	Group(prefix string, fn func(r router.Router)) router.Router

	// With adds one or multiple middleware functions to a specific set of routes.
	// It takes one or more middleware functions as parameters and returns a pointer to the router.Router.
	// This method provides a way to apply middleware to a subset of routes.
	With(middleware ...func(http.Handler) http.Handler) *router.Router
}

func NewRouter() Router {
	return router.NewRouter()
}
