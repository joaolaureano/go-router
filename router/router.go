package router

import (
	"net/http"

	"web/chain"
	"web/const"
	"web/router/context"
	"web/tree"
)

type Router struct {
	root tree.RouterTree

	chain chain.Middleware

	notFound http.HandlerFunc

	closed bool
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
	route := router.root.FindRoute(ctx, _const.HTTPMethods(method), uri)
	if route != nil {
		routeHandler := route.Method[_const.HTTPMethods(r.Method)].Handler
		routeHandler.ServeHTTP(w, r)
	} else {
		router.notFound(w, r)
	}
}

func (router *Router) Register(httpMethod _const.HTTPMethods, path string, method http.HandlerFunc) {
	router.closed = true
	router.root.RegisterRoute(httpMethod,
		path,
		router.chain.BuildHandler(method))
}

func (router *Router) Use(middleware func(http.Handler) http.Handler) {
	if router.closed {
		panic("unable to define middleware after creating first route")
	}
	router.chain.Add(middleware)
}

func (router *Router) NotFound(notFoundFn http.HandlerFunc) {
	router.notFound = notFoundFn
}

func (router *Router) Group(fn func(r Router)) Router {
	chain := chain.NewChain(router.chain.Middlewares()...)
	subrouter := &Router{
		root:     router.root,
		chain:    chain,
		notFound: http.NotFound,
	}
	fn(*subrouter)

	return *subrouter
}
