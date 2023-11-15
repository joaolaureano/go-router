package main

func main() {

	//router := chi.NewRouter()
	//path := "/path"
	//method := func(w http.ResponseWriter, r *http.Request) {
	//
	//	w.Write([]byte(fmt.Sprintf("Hello World Route")))
	//}
	//
	//fn := func(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		r.Header.Set("kkk", "lol")
	//		next.ServeHTTP(w, r)
	//		w.Write([]byte(fmt.Sprintf("Hello World Middleware 1")))
	//	})
	//}
	//
	//router.Get(path, method)
	//router.Use(fn)
	//router.NotFound(func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte(fmt.Sprintf("bug k")))
	//},
	//)

	// Slow handlers/operations.
	//router.Group(func(r chi.Router) {
	//	// Stop processing after 2.5 seconds.
	//	router.Use(func(next http.Handler) http.Handler {
	//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//			r.Header.Set("kkk", "lol")
	//			next.ServeHTTP(w, r)
	//			w.Write([]byte(fmt.Sprintf("Hello World Middleware 2")))
	//		})
	//	})
	//
	//	router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
	//
	//		w.Write([]byte(fmt.Sprintf("/slow method")))
	//	})
	//
	//})
	//
	//http.ListenAndServe("localhost:8080", router)
}
