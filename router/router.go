package router

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/joaolaureano/go-router/chain"
	"github.com/joaolaureano/go-router/router/context"
	"github.com/joaolaureano/go-router/routing"
)

// Method and the verbs below re-export the routing vocabulary, so that naming
// a method does not force callers to import the tree package.
type Method = routing.Method

// WildcardParam is the name a catch-all route captures under: a request to
// "/files/a/b" against "/files/*" reads "a/b" from context.Param(r, "*").
const WildcardParam = routing.WildcardParam

const (
	GET     = routing.GET
	HEAD    = routing.HEAD
	POST    = routing.POST
	PUT     = routing.PUT
	PATCH   = routing.PATCH
	DELETE  = routing.DELETE
	OPTIONS = routing.OPTIONS
)

// routes is the routing tree pinned to what this adapter resolves a route to.
// The tree itself is agnostic; naming the endpoint type is the adapter's job.
type routes = routing.Tree[http.Handler]

type Router struct {
	// state is shared with every router derived from this one, so a route or a
	// fallback handler installed on a group reaches the router that serves.
	state *state

	// chain and prefix belong to this router alone: that is what makes a group
	// a group.
	chain  chain.Middleware
	prefix string
}

func NewRouter() *Router {
	return NewPrefixRouter("")
}

func NewPrefixRouter(prefix string) *Router {
	routeTree := routing.CreateTree[http.Handler]()

	return &Router{
		state:  newState(&routeTree),
		chain:  &chain.Chain{},
		prefix: prefix,
	}
}

// decodeSegment undoes percent-escaping on one path segment.
//
// It is applied per segment rather than to the whole path because URL.Path --
// the path already decoded -- comes too late: a %2F a client escaped precisely
// so it would stay inside one segment has become a separator by then, and no
// route variable could ever hold a slash.
func decodeSegment(segment string) string {
	decoded, err := url.PathUnescape(segment)
	if err != nil {
		// Malformed escaping is not something to guess at; matching the segment
		// literally simply fails to route, which is the honest outcome.
		return segment
	}
	return decoded
}

func defaultMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpMethod := routing.Method(r.Method)

	// Walking the escaped path directly keeps an ordinary request allocation
	// free; only one carrying an escape pays to have its segments decoded.
	path := r.URL.EscapedPath()
	var decode func(string) string
	if strings.IndexByte(path, '%') >= 0 {
		decode = decodeSegment
	}

	// Atomic loads and no lock: see the comment on state.
	router.state.markServed()
	tree := router.state.tree.Load()
	match, status := tree.LookupDecoded(httpMethod, path, decode)
	// RFC 9110: HEAD is GET without content, and net/http already suppresses
	// the body, so a GET route answers HEAD unless one was registered for it.
	if status != routing.StatusFound && httpMethod == HEAD {
		if getMatch, getStatus := tree.LookupDecoded(GET, path, decode); getStatus == routing.StatusFound {
			match, status = getMatch, getStatus
		}
	}

	switch status {
	case routing.StatusFound:
		// A route with no variables has nothing to carry, and deriving a
		// request costs an allocation, so only routes that captured something
		// pay for the context.
		if len(match.Params) > 0 {
			r = context.WithParams(r, match.Params)
		}
		match.Handler.ServeHTTP(w, r)
	case routing.StatusMethodNotAllowed:
		// Allow is set before the handler runs, so a custom one inherits it and
		// can still override it.
		w.Header().Set("Allow", advertisedMethods(match.AllowedMethods))
		if httpMethod == OPTIONS {
			// Nothing was registered for OPTIONS, but the path exists and the
			// header just described it. Answering beats refusing.
			w.WriteHeader(http.StatusNoContent)
			return
		}
		(*router.state.methodNotAllowed.Load())(w, r)
	default:
		(*router.state.notFound.Load())(w, r)
	}
}

// advertisedMethods renders an Allow header from the methods a path was
// registered under, plus the two this router answers on their behalf: OPTIONS,
// which it handles when nothing else does, and HEAD wherever there is a GET.
func advertisedMethods(registered []routing.Method) string {
	advertised := make([]string, 0, len(registered)+2)
	var hasGet, hasHead, hasOptions bool
	for _, httpMethod := range registered {
		switch httpMethod {
		case GET:
			hasGet = true
		case HEAD:
			hasHead = true
		case OPTIONS:
			hasOptions = true
		}
		advertised = append(advertised, string(httpMethod))
	}
	if hasGet && !hasHead {
		advertised = append(advertised, string(HEAD))
	}
	if !hasOptions {
		advertised = append(advertised, string(OPTIONS))
	}
	slices.Sort(advertised)
	return strings.Join(advertised, ", ")
}

func (router *Router) Register(httpMethod Method, path string, method http.HandlerFunc) {
	if method == nil {
		panic(ErrNilHandler)
	}
	if router.prefix != "" {
		path = router.prefix + path
	}
	router.state.mu.Lock()
	defer router.state.mu.Unlock()
	handler := router.chain.BuildHandler(method)
	router.state.mutateTree(func(tree *routes) {
		tree.RegisterRoute(httpMethod, path, handler)
	})
}

// Get and the shortcuts below are Register with the verb spelled into the name,
// which is how routing tables usually read.
func (router *Router) Get(path string, handler http.HandlerFunc) {
	router.Register(GET, path, handler)
}

func (router *Router) Head(path string, handler http.HandlerFunc) {
	router.Register(HEAD, path, handler)
}

func (router *Router) Post(path string, handler http.HandlerFunc) {
	router.Register(POST, path, handler)
}

func (router *Router) Put(path string, handler http.HandlerFunc) {
	router.Register(PUT, path, handler)
}

func (router *Router) Patch(path string, handler http.HandlerFunc) {
	router.Register(PATCH, path, handler)
}

func (router *Router) Delete(path string, handler http.HandlerFunc) {
	router.Register(DELETE, path, handler)
}

func (router *Router) Options(path string, handler http.HandlerFunc) {
	router.Register(OPTIONS, path, handler)
}

func (router *Router) Use(middleware func(http.Handler) http.Handler) {
	router.state.mu.Lock()
	defer router.state.mu.Unlock()
	if router.chain.Sealed() {
		panic(ErrChainSealed)
	}
	router.chain.Add(middleware)
}

func (router *Router) NotFound(notFoundFn http.HandlerFunc) {
	if notFoundFn == nil {
		panic(ErrNilHandler)
	}
	router.state.setNotFound(notFoundFn)
}

// MethodNotAllowed sets the handler for requests whose path exists but not
// under the requested method. The Allow header is already set when it runs.
func (router *Router) MethodNotAllowed(methodNotAllowedFn http.HandlerFunc) {
	if methodNotAllowedFn == nil {
		panic(ErrNilHandler)
	}
	router.state.setMethodNotAllowed(methodNotAllowedFn)
}

// Mount grafts another router's routes in under prefix. The mounted routes keep
// the handlers they were built with, middleware included, so a router assembled
// elsewhere can be attached without knowing anything about this one.
//
// Where both sides define the same method on the same path, this router wins.
func (router *Router) Mount(prefix string, other *Router) {
	if other == nil {
		panic(ErrNilRouter)
	}
	// A group or a With shares the whole state, so mounting one of those would
	// both deadlock on the writer lock and graft the tree into itself.
	if other.state == router.state {
		panic(ErrSharedRoutingTree)
	}

	prefix = router.prefix + prefix
	router.state.mu.Lock()
	defer router.state.mu.Unlock()
	source := other.state.tree.Load()
	router.state.mutateTree(func(tree *routes) {
		tree.MergeAt(prefix, source)
	})
}

// Group returns a subrouter that registers under prefix and starts from a copy
// of this router's middleware. It shares the routing tree, so routes declared
// on it are served by the router this was called on.
func (router *Router) Group(prefix string, fn func(r *Router)) *Router {
	router.state.mu.Lock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	router.state.mu.Unlock()
	subrouter := &Router{
		state:  router.state,
		chain:  chain.NewChain(middlewares...),
		prefix: router.prefix + prefix,
	}

	fn(subrouter)

	return subrouter
}

func (router *Router) With(middleware ...func(http.Handler) http.Handler) *Router {
	router.state.mu.Lock()
	middlewares := append([]func(http.Handler) http.Handler(nil), router.chain.Middlewares()...)
	router.state.mu.Unlock()

	return &Router{
		state:  router.state,
		chain:  chain.NewChain(append(middlewares, middleware...)...),
		prefix: router.prefix,
	}
}
