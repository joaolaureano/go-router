package chain

import "net/http"

type Middleware interface {
	Add(middleware func(handler http.Handler) http.Handler)
	BuildHandler(endpoint http.Handler) http.Handler
	Middlewares() []func(handler http.Handler) http.Handler
	Sealed() bool
}

type Chain struct {
	middlewares []func(http.Handler) http.Handler
	sealed      bool
}

func (chain *Chain) Add(middleware func(handler http.Handler) http.Handler) {
	chain.middlewares = append(chain.middlewares, middleware)
}

// Sealed reports whether the chain has already been baked into a handler.
// Middleware added after that point would apply to later endpoints only,
// silently splitting the chain, so callers must reject it.
func (chain *Chain) Sealed() bool {
	return chain.sealed
}

func (chain *Chain) BuildHandler(endpoint http.Handler) http.Handler {
	chain.sealed = true
	if len(chain.middlewares) == 0 {
		return endpoint
	}
	handler := endpoint
	for i := len(chain.middlewares) - 1; i >= 0; i-- {
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
