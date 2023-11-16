package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_const "web/const"
	"web/router/context"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, router.root, "Root should not be nil")

}
func TestNewRouterWithPrefix(t *testing.T) {
	router := NewPrefixRouter("/prefix")

	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, router.root, "Root should not be nil")
	assert.Equal(t, "/prefix", router.prefix, "Root should not be nil")

}
func TestRouter_RegisterSimplePath(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	r.Register(_const.GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_NotFound(t *testing.T) {
	r := NewRouter()
	path := "/not-found"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	r.NotFound(method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_RegisterWithMiddleware(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware", "HIT")
			next.ServeHTTP(w, r)
		})
	}
	r.Use(middleware)
	r.Register(_const.GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
	assert.Equal(t, "HIT", res.Header.Get("Middleware"))
}
func TestRouter_RegisterBeforeMiddleware(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware", "HIT")
			next.ServeHTTP(w, r)
		})
	}
	r.Register(_const.GET, path, method)

	assert.Panics(t, func() { r.Use(middleware) }, "Use should panic after declaring first route")

}
func TestRouter_RegisterCapturePathVariable(t *testing.T) {
	r := NewRouter()
	path := "/path/{test}"
	pathToTest := "/path/hello"
	method := func(w http.ResponseWriter, r *http.Request) {
		rp := r.Context().Value("RouterContext").(*context.RouterContext)
		w.Write([]byte(rp.Value("test")))
	}
	r.Register(_const.GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, pathToTest))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello", string(body))
}
func TestRouter_RegisterCaptureMultiplePathVariable(t *testing.T) {
	r := NewRouter()
	path := "/path/{test}/path/{test2}"
	pathToTest := "/path/hello/path/world"
	method := func(w http.ResponseWriter, r *http.Request) {
		rp := r.Context().Value("RouterContext").(*context.RouterContext)
		w.Write([]byte(fmt.Sprintf("%s_%s", rp.Value("test"), rp.Value("test2"))))
	}
	r.Register(_const.GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, pathToTest))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_RegisterMultiplePath(t *testing.T) {
	r := NewRouter()
	path1 := "/path/"
	path2 := "/path/test/complex"
	method1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world_1"))
	}
	method2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world_2"))
	}
	r.Register(_const.GET, path1, method1)
	r.Register(_const.GET, path2, method2)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path1))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world_1", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, path2))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, "hello_world_2", string(body))
}
func TestRouter_Group(t *testing.T) {
	router := NewRouter()
	path1 := "/path1"
	group := "/group"
	path2 := "/path2"
	method := func(w http.ResponseWriter, r *http.Request) {
	}
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			w.Write([]byte(fmt.Sprintf("Hello World Middleware 1")))
		})
	}
	router.Use(fn)
	router.Register(_const.GET, path1, method)
	router.Group(group, func(r Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				w.Write([]byte(fmt.Sprintf("Hello World Middleware 2")))
			})
		})
		r.Register(_const.GET, path2, method)
	})
	s := setup(router)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path1))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Hello World Middleware 1", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, group+path2))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, "Hello World Middleware 1Hello World Middleware 2", string(body))
}
func setup(r http.Handler) *httptest.Server {

	return httptest.NewServer(r)
}
