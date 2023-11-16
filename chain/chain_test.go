package chain

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain_SingleMiddleware(t *testing.T) {
	// Create a new Chain
	c := &Chain{}

	// Define a sample handler
	sampleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add middleware to the chain
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware", "HIT")
			next.ServeHTTP(w, r)
		})
	}

	c.Add(middleware)

	// Build the handler with the chain
	handler := c.BuildHandler(sampleHandler)

	// Create a test server with the built handler
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Make a request to the test server
	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer res.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "HIT", res.Header.Get("Middleware"))
}

func TestChain_MultipleMiddleware(t *testing.T) {
	// Create a new Chain
	c := &Chain{}
	expectedOutput := "hello world 1hello world 2hello world 3"

	// Define a sample handler
	sampleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add middleware to the chain
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware-1", "HIT-1")
			next.ServeHTTP(w, r)
			w.Write([]byte("hello world 1"))
		})
	}
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware-2", "HIT-2")
			next.ServeHTTP(w, r)
			w.Write([]byte("hello world 2"))
		})
	}
	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware-3", "HIT-3")
			next.ServeHTTP(w, r)
			w.Write([]byte("hello world 3"))
		})
	}

	c.Add(middleware1)
	c.Add(middleware2)
	c.Add(middleware3)

	// Build the handler with the chain
	handler := c.BuildHandler(sampleHandler)

	// Create a test server with the built handler
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Make a request to the test server
	res, err := http.Get(ts.URL)
	defer res.Body.Close()

	// Assert the response status code
	body, _ := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expectedOutput, string(body))
	assert.Equal(t, "HIT-1", res.Header.Get("Middleware-1"))
	assert.Equal(t, "HIT-2", res.Header.Get("Middleware-2"))
	assert.Equal(t, "HIT-3", res.Header.Get("Middleware-3"))
}

func TestChain_NoMiddleware(t *testing.T) {
	// Create a new Chain
	c := &Chain{}

	// Define a sample handler
	sampleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Build the handler with the chain
	handler := c.BuildHandler(sampleHandler)

	// Create a test server with the built handler
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Make a request to the test server
	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer res.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestNewChain(t *testing.T) {
	chain := NewChain()

	assert.NotNil(t, chain, "Chain should not be nil")

}

func TestNewChain_Middlewares(t *testing.T) {
	chain := NewChain()

	assert.NotNil(t, len(chain.Middlewares()), "Chain should not be nil")

}
