package main

import (
	"net/http"

	"web/router"
	router2 "web/router/tree/router"
)

// Router is an interface that extends the http.Handler interface,
// requiring any implementing type to fulfill the http.Handler contract.
// It declares two methods - Register and Use - for registering routes and adding middleware.
type Router interface {
	http.Handler

	// Register is a method for adding a new route to the router.
	// It takes an HTTP method, a path, and a handler function as parameters.
	// The router will use these parameters to associate incoming requests with the specified handler.
	Register(httpMethod router.HTTPMethods, path string, method http.HandlerFunc)

	// Use is a method for adding middleware to the router.
	// Middleware functions can process or modify requests before reaching the route handler.
	// This method enhances the router's functionality by allowing the insertion of additional processing steps.
	Use(middleware func(http.Handler) http.Handler)
}

func NewRouter() Router {
	return router2.NewRouter()
}