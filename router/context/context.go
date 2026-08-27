package context

import (
	"context"
	"net/http"
)

const RouterContextKey = "RouterContext"

type contextKey struct{}

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

func (routerCtx *RouterContext) InjectIntoRequest(r *http.Request) {
	*r = *routerCtx.WithRequest(r)
}

func (routerCtx *RouterContext) WithRequest(r *http.Request) *http.Request {
	requestContext := context.WithValue(r.Context(), contextKey{}, routerCtx)
	return r.WithContext(context.WithValue(requestContext, RouterContextKey, routerCtx))
}

func FromRequest(r *http.Request) (*RouterContext, bool) {
	routerCtx, ok := r.Context().Value(contextKey{}).(*RouterContext)
	if !ok {
		routerCtx, ok = r.Context().Value(RouterContextKey).(*RouterContext)
	}
	return routerCtx, ok
}
