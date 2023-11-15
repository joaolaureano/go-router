package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	router := chi.NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(fmt.Sprintf("Hello World Route")))
	}

	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("kkk", "lol")
			next.ServeHTTP(w, r)
			w.Write([]byte(fmt.Sprintf("Hello World Middleware")))
		})
	}

	router.Use(fn)
	router.Get(path, method)
	//router.NotFound(func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte(fmt.Sprintf("bug k")))
	//},
	//)

	// Slow handlers/operations.
	//router.Group(func(r chi.Router) {
	//	// Stop processing after 2.5 seconds.
	//	r.Use(func(next http.Handler) http.Handler {
	//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//			r.Header.Set("kkk", "lol")
	//			next.ServeHTTP(w, r)
	//			w.Write([]byte(fmt.Sprintf("group1")))
	//		})
	//	})
	//
	//	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
	//		rand.Seed(time.Now().Unix())
	//
	//		// Processing will take 1-5 seconds.
	//		processTime := time.Duration(rand.Intn(4)+1) * time.Second
	//
	//		select {
	//		case <-r.Context().Done():
	//			return
	//
	//		case <-time.After(processTime):
	//			// The above channel simulates some hard work.
	//		}
	//
	//		w.Write([]byte(fmt.Sprintf("/slow method")))
	//	})
	//
	//	router.Group(func(r chi.Router) {
	//		// Stop processing after 2.5 seconds.
	//		r.Use(func(next http.Handler) http.Handler {
	//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//				r.Header.Set("kkk", "lol")
	//				next.ServeHTTP(w, r)
	//				w.Write([]byte(fmt.Sprintf("group2")))
	//			})
	//		})
	//
	//		r.Get("/slow2", func(w http.ResponseWriter, r *http.Request) {
	//			rand.Seed(time.Now().Unix())
	//
	//			// Processing will take 1-5 seconds.
	//			processTime := time.Duration(rand.Intn(4)+1) * time.Second
	//
	//			select {
	//			case <-r.Context().Done():
	//				return
	//
	//			case <-time.After(processTime):
	//				// The above channel simulates some hard work.
	//			}
	//
	//			w.Write([]byte(fmt.Sprintf("/slow2 method")))
	//		})
	//	})
	//})

	http.ListenAndServe("localhost:8080", router)
}
