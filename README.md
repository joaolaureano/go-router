[![en](https://img.shields.io/badge/lang-en-blue.svg)](https://github.com/joaolaureano/gorouter/blob/main/README.md)
[![pt-br](https://img.shields.io/badge/lang-pt--br-green.svg)](https://github.com/joaolaureano/gorouter/blob/main/README.pt-PT.md)

# Go-Router

Go-Router is a simple Golang router to handle HTTP requests.
The project was developed for the purpose of learning and experimenting with the Go language.


## Install
```go get github.com/joaolaureano/go-router@latest```

### Code example
**As easy as**
```go
    package main
    
    import (
        "fmt"
        "net/http"
        
        "github.com/joaolaureano/go-router/router"
        "github.com/joaolaureano/go-router/router/context"
    )
    
    func main() {
    r := router.NewRouter()
    
        r.Register(router.GET, "/ping", func(writer http.ResponseWriter, request *http.Request) {
            writer.Write([]byte("pong"))
        })
        r.Group("/{id}", func(r *router.Router) {
            r.Use(func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
                    next.ServeHTTP(writer, request)
                    message := fmt.Sprintf("Group Middleware\n Found Path value: %s", context.Param(request, "id"))
                    writer.Write([]byte(message))
                })
            })
            r.Register(router.GET, "/pong", func(writer http.ResponseWriter, request *http.Request) {
                writer.Write([]byte("ping"))
            })
        })
        http.ListenAndServe(":3333", r)
    }
```
You can find more at folder ```.example/```


### Catch-all routes
A trailing `*` segment matches the rest of the path and captures it under `router.WildcardParam`:

```go
r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, filepath.Join("./public", context.Param(r, router.WildcardParam)))
})
```

Static and parameter routes take precedence, so `/files/readme` still reaches a route registered for it. `*` is only a catch-all as a whole segment, and only as the last one.

## Interface
- `Register(httpMethod router.Method, path string, method http.HandlerFunc)`: Registers an HTTP method for a specific path.
- `Get/Head/Post/Put/Patch/Delete/Options(path string, handler http.HandlerFunc)`: Shortcuts for `Register` with the verb spelled into the name.
- `Use(middleware func(http.Handler) http.Handler)`: Uses middleware to handle HTTP requests.
- `NotFound(notFoundFn http.HandlerFunc)`: Sets a handler for requests on non-existent routes.
- `MethodNotAllowed(fn http.HandlerFunc)`: Sets a handler for requests to an existing path under an unregistered method.
- `Group(prefix string, fn func(r *router.Router)) *router.Router`: Groups routes under a specified prefix.
- `Mount(prefix string, other *router.Router)`: Grafts another router's routes under a prefix, handlers and middleware intact.
- `With(middleware ...func(http.Handler) http.Handler) *router.Router`: Uses middleware for a specific set of routes.

## Credits

This project was inspired and influenced by **[go-Chi](https://github.com/go-chi/chi)**.


## Contributing

Feel free to open issues or send pull requests to contribute to improvements in this project.
Every contribution is welcomed!
