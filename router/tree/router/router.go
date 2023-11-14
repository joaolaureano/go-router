package router

import (
	"net/http"

	router2 "web/router"
	"web/router/context"
	"web/router/tree"
	"web/router/tree/chain"
)

type Router struct {
	root tree.RouterTree

	chain chain.Middleware

	notFound http.HandlerFunc
}

func NewRouter() *Router {
	tree := tree.CreateTree()

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: http.NotFound,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	method := r.Method
	ctx := context.NewContext()
	ctx.InjectIntoRequest(r)
	route := router.root.FindRoute(ctx, router2.HTTPMethods(method), uri)
	if route != nil {
		routeHandler := route.Method[router2.HTTPMethods(r.Method)].Handler
		router.chain.BuildHandler(routeHandler).ServeHTTP(w, r)
	} else {
		router.notFound(w, r)
	}
}

func (router *Router) Register(httpMethod router2.HTTPMethods, path string, method http.HandlerFunc) {
	router.root.RegisterRoute(httpMethod, path, method)
}

func (router *Router) Use(middleware func(http.Handler) http.Handler) {
	router.chain.Add(middleware)
}

func (router *Router) NotFound(notFoundFn http.HandlerFunc) {
	router.notFound = notFoundFn
}
