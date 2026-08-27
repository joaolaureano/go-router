package router

import (
	"net/http"
	"sync"

	"github.com/joaolaureano/go-router/chain"
	_const "github.com/joaolaureano/go-router/const"
	"github.com/joaolaureano/go-router/router/context"
	"github.com/joaolaureano/go-router/tree"
)

type Router struct {
	root tree.RouterTree

	chain chain.Middleware

	notFound http.HandlerFunc

	prefix string

	mu *sync.RWMutex
}

func NewRouter() *Router {
	tree := tree.CreateTree()

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: http.NotFound,
		mu:       &sync.RWMutex{},
	}
}

func NewPrefixRouter(prefix string) *Router {
	tree := tree.CreateTree()

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: http.NotFound,
		prefix:   prefix,
		mu:       &sync.RWMutex{},
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	method := r.Method
	ctx := context.NewContext()
	r = ctx.WithRequest(r)
	router.mu.RLock()
	route := router.root.FindRoute(ctx, _const.HTTPMethods(method), uri)
	var routeHandler http.Handler
	pathExists := false
	if route != nil {
		routeHandler = route.Method[_const.HTTPMethods(r.Method)].Handler
	} else {
		pathExists = router.root.FindPath(uri) != nil
	}
	notFound := router.notFound
	router.mu.RUnlock()
	if routeHandler != nil {
		routeHandler.ServeHTTP(w, r)
	} else if pathExists {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	} else {
		notFound(w, r)
	}
}

func (router *Router) Register(httpMethod _const.HTTPMethods, path string, method http.HandlerFunc) {
	if method == nil {
		panic("handler must not be nil")
	}
	router.mu.Lock()
	defer router.mu.Unlock()
	if router.prefix != "" {
		path = router.prefix + path
	}
	router.root.RegisterRoute(httpMethod,
		path,
		router.chain.BuildHandler(method))
}

func (router *Router) Use(middleware func(http.Handler) http.Handler) {
	router.mu.Lock()
	defer router.mu.Unlock()
	if router.chain.Sealed() {
		panic("unable to define middleware after creating first route")
	}
	router.chain.Add(middleware)
}

func (router *Router) NotFound(notFoundFn http.HandlerFunc) {
	router.mu.Lock()
	defer router.mu.Unlock()
	router.notFound = notFoundFn
}

func (router *Router) Group(prefix string, fn func(r Router)) Router {
	router.mu.RLock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	fullPrefix := router.prefix + prefix
	notFound := router.notFound
	mu := router.mu
	router.mu.RUnlock()
	chain := chain.NewChain(middlewares...)
	subrouter := &Router{
		root:     router.root,
		chain:    chain,
		notFound: notFound,
		prefix:   fullPrefix,
		mu:       mu,
	}

	fn(*subrouter)

	return *subrouter
}

func (router *Router) With(middleware ...func(http.Handler) http.Handler) *Router {
	router.mu.RLock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	subrouter := &Router{
		root:     router.root,
		notFound: router.notFound,
		prefix:   router.prefix,
		mu:       router.mu,
	}
	router.mu.RUnlock()
	subrouter.chain = chain.NewChain(append(middlewares, middleware...)...)

	return subrouter
}
