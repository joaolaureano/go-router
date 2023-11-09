package main

import (
	"fmt"
	"net/http"

	router2 "web/router"
	"web/router/tree/router"
)

func main() {

	router := router.NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(fmt.Sprintf("Hello World Route")))
	}

	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("kkk", "lol")
			w.Write([]byte(fmt.Sprintf("Hello World Middleware")))
			next.ServeHTTP(w, r)
		})
	}

	router.Use(fn)
	router.Register(router2.GET, path, method)

	http.ListenAndServe("localhost:8080", router)
}
