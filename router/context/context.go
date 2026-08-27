package context

import (
	"context"
	"net/http"
)

// contextKey is unexported and of a package-local type, so no other package can
// construct it. That is the point: a context key of a basic type such as string
// collides silently with any other package that picks the same text.
type contextKey struct{}

// RouterContext holds the route variables captured for one request.
//
// A slice rather than a map: a route declares a handful of variables at most,
// and at that size a linear scan beats hashing while costing one allocation
// instead of a map header plus buckets.
type RouterContext struct {
	params []param
}

type param struct {
	key   string
	value string
}

// Value reads a route variable. It tolerates a nil receiver so that the common
// `routerCtx, _ := FromRequest(r)` shape stays safe on a route that captured
// nothing and therefore carries no context.
func (routerCtx *RouterContext) Value(key string) string {
	if routerCtx == nil {
		return ""
	}
	for _, entry := range routerCtx.params {
		if entry.key == key {
			return entry.value
		}
	}
	return ""
}

func (routerCtx *RouterContext) Set(key string, value string) {
	for i := range routerCtx.params {
		if routerCtx.params[i].key == key {
			routerCtx.params[i].value = value
			return
		}
	}
	if routerCtx.params == nil {
		// Room for a typical route's variables up front, so filling the context
		// does not regrow the slice once per variable.
		routerCtx.params = make([]param, 0, 4)
	}
	routerCtx.params = append(routerCtx.params, param{key: key, value: value})
}

func NewContext() *RouterContext {
	return &RouterContext{}
}

// WithRequest returns a copy of r carrying this context.
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
