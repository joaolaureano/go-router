package context

import (
	"context"
	"net/http"

	"github.com/joaolaureano/go-router/routing"
)

// contextKey is unexported and of a package-local type, so no other package can
// construct it. That is the point: a context key of a basic type such as string
// collides silently with any other package that picks the same text.
type contextKey struct{}

// RouterContext holds the route variables captured for one request.
//
// A slice rather than a map: a route declares a handful of variables at most,
// and at that size a linear scan beats hashing while costing one allocation
// instead of a map header plus buckets. It is the tree's own Param slice, which
// a lookup builds fresh each time and hands over rather than have the router
// copy it into a second one.
type RouterContext struct {
	params []routing.Param
}

// Value reads a route variable. It tolerates a nil receiver so that the common
// `routerCtx, _ := FromRequest(r)` shape stays safe on a route that captured
// nothing and therefore carries no context.
func (routerCtx *RouterContext) Value(key string) string {
	if routerCtx == nil {
		return ""
	}
	for _, entry := range routerCtx.params {
		if entry.Name == key {
			return entry.Value
		}
	}
	return ""
}

func (routerCtx *RouterContext) Set(key string, value string) {
	for i := range routerCtx.params {
		if routerCtx.params[i].Name == key {
			routerCtx.params[i].Value = value
			return
		}
	}
	if routerCtx.params == nil {
		// Room for a typical route's variables up front, so filling the context
		// does not regrow the slice once per variable.
		routerCtx.params = make([]routing.Param, 0, 4)
	}
	routerCtx.params = append(routerCtx.params, routing.Param{Name: key, Value: value})
}

func NewContext() *RouterContext {
	return &RouterContext{}
}

// FromParams takes ownership of the params a lookup produced. A lookup builds
// that slice fresh for each request, so there is nothing to copy and nothing
// shared with the routing table.
func FromParams(params []routing.Param) *RouterContext {
	return &RouterContext{params: params}
}

// routeContext carries the route variables by being a context rather than by
// sitting inside one. context.WithValue allocates a node holding a pointer to a
// second allocation; being the node folds the two together, and the pointer
// handed back to FromRequest points inside this same object.
type routeContext struct {
	context.Context
	routerCtx RouterContext
}

func (c *routeContext) Value(key any) any {
	if _, ours := key.(contextKey); ours {
		return &c.routerCtx
	}
	return c.Context.Value(key)
}

// WithParams returns a copy of r carrying these route variables, in a single
// allocation.
func WithParams(r *http.Request, params []routing.Param) *http.Request {
	return r.WithContext(&routeContext{
		Context:   r.Context(),
		routerCtx: RouterContext{params: params},
	})
}

// WithRequest returns a copy of r carrying this context. It attaches this very
// object, not a copy, so that a Set afterwards is still visible through the
// request -- which is why it does not use the folded form above.
func (routerCtx *RouterContext) WithRequest(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextKey{}, routerCtx))
}

// FromRequest returns the routing context attached to r, if there is one.
func FromRequest(r *http.Request) (*RouterContext, bool) {
	routerCtx, ok := r.Context().Value(contextKey{}).(*RouterContext)
	return routerCtx, ok
}

// Param returns a single route variable, or the empty string when the request
// carries no routing context or the route declares no such variable.
func Param(r *http.Request, key string) string {
	routerCtx, ok := FromRequest(r)
	if !ok {
		return ""
	}
	return routerCtx.Value(key)
}
