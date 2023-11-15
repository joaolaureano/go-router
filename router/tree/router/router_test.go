package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"web/router"
	"web/router/context"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, router.root, "Root should not be nil")

}
func TestRouter_RegisterSimplePath(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	r.Register(router.GET, path, method)
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

	r.Register(router.GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
	assert.Equal(t, "HIT", res.Header.Get("Middleware"))
}
func TestRouter_RegisterCapturePathVariable(t *testing.T) {
	r := NewRouter()
	path := "/path/{test}"
	pathToTest := "/path/hello"
	method := func(w http.ResponseWriter, r *http.Request) {
		rp := r.Context().Value("RouterContext").(*context.RouterContext)
		w.Write([]byte(rp.Value("test")))
	}
	r.Register(router.GET, path, method)
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
	r.Register(router.GET, path, method)
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
	r.Register(router.GET, path1, method1)
	r.Register(router.GET, path2, method2)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path1))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world_1", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, path2))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, "hello_world_2", string(body))
}
func setup(r http.Handler) *httptest.Server {

	return httptest.NewServer(r)
}
