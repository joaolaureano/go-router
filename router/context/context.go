package context

import (
	"context"
	"net/http"
	"slices"
)

const RouterContextKey = "RouterContext"

type RouterContext struct {
	paramNames []string
	paramValue []string
}

func (routerCtx *RouterContext) Value(key string) string {
	idx := slices.Index(routerCtx.paramNames, key)
	if idx > -1 {
		return routerCtx.paramValue[idx]
	}
	return ""
}

func (routerCtx *RouterContext) Set(key string, value string) {
	routerCtx.paramNames = append(routerCtx.paramNames, key)
	routerCtx.paramValue = append(routerCtx.paramValue, value)
}

func NewContext() *RouterContext {
	return &RouterContext{
		paramNames: make([]string, 0),
		paramValue: make([]string, 0),
	}
}

func (routerCtx *RouterContext) InjectIntoRequest(r *http.Request) {
	*r = *r.WithContext(context.WithValue((*r).Context(), RouterContextKey, routerCtx))
}
