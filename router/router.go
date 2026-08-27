package router

import (
	"net/http"

	"strings"
	"sync"

	"github.com/joaolaureano/go-router/chain"
	_const "github.com/joaolaureano/go-router/const"
	"github.com/joaolaureano/go-router/router/context"
	"github.com/joaolaureano/go-router/tree"
)

type Router struct {
	root tree.RouterTree

	chain chain.Middleware

	// notFound is shared by reference across the router family, like root and
	// mu, so a handler installed on a group reaches the router that serves.
	notFound *http.HandlerFunc

	prefix string

	mu *sync.RWMutex
}

func NewRouter() *Router {
	tree := tree.CreateTree()
	notFound := http.HandlerFunc(http.NotFound)

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: &notFound,
		mu:       &sync.RWMutex{},
	}
}

func NewPrefixRouter(prefix string) *Router {
	tree := tree.CreateTree()
	notFound := http.HandlerFunc(http.NotFound)

	return &Router{
		root:     &tree,
		chain:    &chain.Chain{},
		notFound: &notFound,
		prefix:   prefix,
		mu:       &sync.RWMutex{},
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routerCtx := context.NewContext()
	r = routerCtx.WithRequest(r)

	router.mu.RLock()
	match, status := router.root.Lookup(_const.HTTPMethods(r.Method), r.URL.Path)
	notFound := *router.notFound
	router.mu.RUnlock()

	switch status {
	case tree.StatusFound:
		for _, param := range match.Params {
			routerCtx.Set(param.Name, param.Value)
		}
		match.Handler.ServeHTTP(w, r)
	case tree.StatusMethodNotAllowed:
		w.Header().Set("Allow", strings.Join(match.AllowedMethods, ", "))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	default:
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
	if notFoundFn == nil {
		panic("handler must not be nil")
	}
	router.mu.Lock()
	defer router.mu.Unlock()
	*router.notFound = notFoundFn
}

func (router *Router) Group(prefix string, fn func(r Router)) Router {
	router.mu.RLock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	fullPrefix := router.prefix + prefix
	mu := router.mu
	router.mu.RUnlock()
	chain := chain.NewChain(middlewares...)
	subrouter := &Router{
		root:     router.root,
		chain:    chain,
		notFound: router.notFound,
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
