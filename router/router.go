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

	prefix string
}

func NewRouter() *Router {
	tree := tree.CreateTree()

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: http.NotFound,
	}
}

func NewPrefixRouter(prefix string) *Router {
	tree := tree.CreateTree()

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: http.NotFound,
		prefix:   prefix,
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
	if router.prefix != "" {
		path = router.prefix + path
	}
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

func (router *Router) Group(prefix string, fn func(r Router)) Router {
	chain := chain.NewChain(router.chain.Middlewares()...)
	//tree := tree.CreateTree()
	subrouter := &Router{
		root:     router.root,
		chain:    chain,
		notFound: http.NotFound,
		prefix:   prefix,
	}

	fn(*subrouter)

	//router.root.Merge(subrouter.root)

	return *subrouter
}

func (router *Router) With(middleware ...func(http.Handler) http.Handler) *Router {
	r := NewRouter()
	r.root = router.root

	for _, m := range middleware {
		r.Use(m)
	}

	return r
}
