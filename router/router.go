package router

import (
	"context"
	"net/http"

	context2 "web/router/context"
)

type Router struct {
	root Tree
}

func NewRouter() *Router {
	tree := createTree()
	return &Router{
		root: tree,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	method := r.Method
	ctx := context2.NewContext()
	r = r.WithContext(context.WithValue(r.Context(), context2.ContextKey, ctx))
	route := router.root.FindRoute(ctx, HTTPMethods(method), uri)
	if route != nil {
		route.method[HTTPMethods(r.Method)].handler.ServeHTTP(w, r)
	}
}

func (router *Router) Register(httpMethod HTTPMethods, path string, method http.HandlerFunc) {
	router.root.RegisterRoute(httpMethod, path, method)
}
