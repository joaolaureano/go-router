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
type RouterContext struct {
	params map[string]string
}

func (routerCtx *RouterContext) Value(key string) string {
	return routerCtx.params[key]
}

func (routerCtx *RouterContext) Set(key string, value string) {
	if routerCtx.params == nil {
		routerCtx.params = make(map[string]string)
	}
	routerCtx.params[key] = value
}

// NewContext leaves params nil. Reads on a nil map are legal and Set allocates
// on demand, so a route without variables never pays for the map.
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
