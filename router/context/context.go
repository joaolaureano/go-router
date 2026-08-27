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

func NewContext() *RouterContext {
	return &RouterContext{
		params: make(map[string]string),
	}
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
