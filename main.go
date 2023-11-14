package main

import (
	"fmt"
	"net/http"

	router2 "web/router"
)

func main() {

	router := NewRouter()
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
	router.Register(router2.GET, path, method)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("bug k")))
	},
	)

	http.ListenAndServe("localhost:8080", router)
}
