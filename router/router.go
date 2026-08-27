package router

import (
	"net/http"

	"strings"
	"sync"

	"github.com/joaolaureano/go-router/chain"
	"github.com/joaolaureano/go-router/router/context"
	"github.com/joaolaureano/go-router/tree"
)

// Method and the verbs below re-export the routing vocabulary, so that naming
// a method does not force callers to import the tree package.
type Method = tree.Method

const (
	GET     = tree.GET
	HEAD    = tree.HEAD
	POST    = tree.POST
	PUT     = tree.PUT
	PATCH   = tree.PATCH
	DELETE  = tree.DELETE
	OPTIONS = tree.OPTIONS
)

// routeTree is the slice of the tree the router actually needs. Declaring it
// here rather than beside Tree keeps the domain free to grow methods without
// widening what the router is coupled to.
type routeTree interface {
	RegisterRoute(httpMethod Method, newValue string, method http.Handler)
	Lookup(httpMethod Method, path string) (tree.Match, tree.Status)
}

type Router struct {
	root routeTree

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
	match, status := router.root.Lookup(tree.Method(r.Method), r.URL.Path)
	notFound := *router.notFound
	router.mu.RUnlock()

	switch status {
	case tree.StatusFound:
		for _, param := range match.Params {
			routerCtx.Set(param.Name, param.Value)
		}
		match.Handler.ServeHTTP(w, r)
	case tree.StatusMethodNotAllowed:
		allowed := make([]string, len(match.AllowedMethods))
		for i, httpMethod := range match.AllowedMethods {
			allowed[i] = string(httpMethod)
		}
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	default:
		notFound(w, r)
	}
}

func (router *Router) Register(httpMethod Method, path string, method http.HandlerFunc) {
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
