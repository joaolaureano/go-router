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

	// notFound and methodNotAllowed are shared by reference across the router
	// family, like root and mu, so a handler installed on a group reaches the
	// router that serves.
	notFound         *http.HandlerFunc
	methodNotAllowed *http.HandlerFunc

	prefix string

	mu *sync.RWMutex
}

func NewRouter() *Router {
	return NewPrefixRouter("")
}

func NewPrefixRouter(prefix string) *Router {
	routes := tree.CreateTree()
	notFound := http.HandlerFunc(http.NotFound)
	methodNotAllowed := http.HandlerFunc(defaultMethodNotAllowed)

	return &Router{
		root:             &routes,
		chain:            &chain.Chain{},
		notFound:         &notFound,
		methodNotAllowed: &methodNotAllowed,
		prefix:           prefix,
		mu:               &sync.RWMutex{},
	}
}

func defaultMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mu.RLock()
	match, status := router.root.Lookup(tree.Method(r.Method), r.URL.Path)
	notFound := *router.notFound
	methodNotAllowed := *router.methodNotAllowed
	router.mu.RUnlock()

	switch status {
	case tree.StatusFound:
		// A route with no variables has nothing to carry, and deriving a
		// request costs an allocation, so only routes that captured something
		// pay for the context.
		if len(match.Params) > 0 {
			routerCtx := context.NewContext()
			for _, param := range match.Params {
				routerCtx.Set(param.Name, param.Value)
			}
			r = routerCtx.WithRequest(r)
		}
		match.Handler.ServeHTTP(w, r)
	case tree.StatusMethodNotAllowed:
		allowed := make([]string, len(match.AllowedMethods))
		for i, httpMethod := range match.AllowedMethods {
			allowed[i] = string(httpMethod)
		}
		// Allow is set before the handler runs, so a custom one inherits it and
		// can still override it.
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		methodNotAllowed(w, r)
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

// MethodNotAllowed sets the handler for requests whose path exists but not
// under the requested method. The Allow header is already set when it runs.
func (router *Router) MethodNotAllowed(methodNotAllowedFn http.HandlerFunc) {
	if methodNotAllowedFn == nil {
		panic("handler must not be nil")
	}
	router.mu.Lock()
	defer router.mu.Unlock()
	*router.methodNotAllowed = methodNotAllowedFn
}

// Group returns a subrouter that registers under prefix and starts from a copy
// of this router's middleware. It shares the routing tree, so routes declared
// on it are served by the router this was called on.
func (router *Router) Group(prefix string, fn func(r *Router)) *Router {
	router.mu.RLock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	fullPrefix := router.prefix + prefix
	mu := router.mu
	router.mu.RUnlock()
	chain := chain.NewChain(middlewares...)
	subrouter := &Router{
		root:             router.root,
		chain:            chain,
		notFound:         router.notFound,
		methodNotAllowed: router.methodNotAllowed,
		prefix:           fullPrefix,
		mu:               mu,
	}

	fn(subrouter)

	return subrouter
}

func (router *Router) With(middleware ...func(http.Handler) http.Handler) *Router {
	router.mu.RLock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	subrouter := &Router{
		root:             router.root,
		notFound:         router.notFound,
		methodNotAllowed: router.methodNotAllowed,
		prefix:           router.prefix,
		mu:               router.mu,
	}
	router.mu.RUnlock()
	subrouter.chain = chain.NewChain(append(middlewares, middleware...)...)

	return subrouter
}
