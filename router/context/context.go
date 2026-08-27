package context

import (
	"context"
	"net/http"
)

const RouterContextKey = "RouterContext"

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
	*r = *r.WithContext(context.WithValue((*r).Context(), RouterContextKey, routerCtx))
}
