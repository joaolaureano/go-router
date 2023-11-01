package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"web/router/context"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, router.root, "Root should not be nil")

}

func TestRouter_RegisterSimplePath(t *testing.T) {
	router := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	router.Register(GET, path, method)
	s := setup(router)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_RegisterCapturePathVariable(t *testing.T) {
	router := NewRouter()
	path := "/path/{test}"
	pathToTest := "/path/hello"
	method := func(w http.ResponseWriter, r *http.Request) {
		rp := r.Context().Value("RouterContext").(*context.Context)
		w.Write([]byte(rp.Value("test")))
	}
	router.Register(GET, path, method)
	s := setup(router)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, pathToTest))

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "hello", string(body))
}
func TestRouter_RegisterCaptureMultiplePathVariable(t *testing.T) {
	router := NewRouter()
	path := "/path/{test}/path/{test2}"
	pathToTest := "/path/hello/path/world"
	method := func(w http.ResponseWriter, r *http.Request) {
		rp := r.Context().Value("RouterContext").(*context.Context)
		w.Write([]byte(fmt.Sprintf("%s_%s", rp.Value("test"), rp.Value("test2"))))
	}
	router.Register(GET, path, method)
	s := setup(router)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, pathToTest))

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_RegisterMultiplePath(t *testing.T) {
	router := NewRouter()
	path1 := "/path/"
	path2 := "/path/test/complex"
	method1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world_1"))
	}
	method2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world_2"))
	}
	router.Register(GET, path1, method1)
	router.Register(GET, path2, method2)
	s := setup(router)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path1))
	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "hello_world_1", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, path2))
	body, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, "hello_world_2", string(body))
}
func setup(r http.Handler) *httptest.Server {

	return httptest.NewServer(r)
}
