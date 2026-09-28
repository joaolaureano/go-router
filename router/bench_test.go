package router

// The benchmarks the routing decisions in this repository were made against.
// They serve one router covering the shapes a real table mixes -- static
// segments, variables, a level fanning out fifty ways -- so that a change
// cannot look good on one shape while quietly taxing another.
//
// Read them with benchstat and -count=10 or more. A single sample has more
// than once suggested a regression that turned out to be a double-digit gain.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func benchRouter() *Router {
	r := NewRouter()
	h := func(w http.ResponseWriter, r *http.Request) {}
	r.Get("/", h)
	r.Get("/users", h)
	r.Get("/users/{id}", h)
	r.Get("/users/{id}/posts/{postID}", h)
	r.Post("/users/{id}", h)
	r.Put("/users/{id}", h)
	r.Delete("/users/{id}", h)
	r.Get("/a/b/c/d/e", h)
	for i := range 50 {
		r.Get(fmt.Sprintf("/res%d", i), h)
	}
	return r
}

func run(b *testing.B, r *Router, method, target string) {
	b.ReportAllocs()
	req := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	b.ResetTimer()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkStatic(b *testing.B)   { run(b, benchRouter(), "GET", "/a/b/c/d/e") }
func BenchmarkParam(b *testing.B)    { run(b, benchRouter(), "GET", "/users/42/posts/7") }
func BenchmarkMiss(b *testing.B)     { run(b, benchRouter(), "GET", "/nope/nothing") }
func BenchmarkEscaped(b *testing.B)  { run(b, benchRouter(), "GET", "/users/a%2Fb") }
func BenchmarkFanout50(b *testing.B) { run(b, benchRouter(), "GET", "/res49") }

func BenchmarkMethodNotAllowed(b *testing.B) { run(b, benchRouter(), "PATCH", "/users/42") }
func BenchmarkOptions(b *testing.B)          { run(b, benchRouter(), "OPTIONS", "/users/42") }

func BenchmarkStaticParallel(b *testing.B) {
	r := benchRouter()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/a/b/c/d/e", nil)
		w := httptest.NewRecorder()
		for pb.Next() {
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkRegister300(b *testing.B) {
	paths := make([]string, 300)
	for i := range paths {
		paths[i] = fmt.Sprintf("/g%d/{id}/sub%d", i%20, i)
	}
	h := func(w http.ResponseWriter, r *http.Request) {}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		r := NewRouter()
		for _, p := range paths {
			r.Get(p, h)
		}
	}
}
