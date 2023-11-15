package chain

import "net/http"

type Middleware interface {
	Add(middleware func(handler http.Handler) http.Handler)
	BuildHandler(endpoint http.Handler) http.Handler
	Middlewares() []func(handler http.Handler) http.Handler
}

type Chain struct {
	middlewares []func(http.Handler) http.Handler
}

func (chain *Chain) Add(middleware func(handler http.Handler) http.Handler) {
	chain.middlewares = append(chain.middlewares, middleware)
}

func (chain *Chain) BuildHandler(endpoint http.Handler) http.Handler {
	if len(chain.middlewares) == 0 {
		return endpoint
	}
	handler := endpoint
	for i := 0; i < len(chain.middlewares); i++ {
		handler = chain.middlewares[i](handler)
	}

	return handler
}

func NewChain(middlewares ...func(http.Handler) http.Handler) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

func (chain *Chain) Middlewares() []func(http.Handler) http.Handler {
	return chain.middlewares
}
