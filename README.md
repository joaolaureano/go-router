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

## Interface
- `Register(httpMethod router.Method, path string, method http.HandlerFunc)`: Registers an HTTP method for a specific path.
- `Use(middleware func(http.Handler) http.Handler)`: Uses middleware to handle HTTP requests.
- `NotFound(notFoundFn http.HandlerFunc)`: Sets a handler for requests on non-existent routes.
- `Group(prefix string, fn func(r *router.Router)) *router.Router`: Groups routes under a specified prefix.
- `With(middleware ...func(http.Handler) http.Handler) *router.Router`: Uses middleware for a specific set of routes.

## Credits

This project was inspired and influenced by **[go-Chi](https://github.com/go-chi/chi)**.


## Contributing

Feel free to open issues or send pull requests to contribute to improvements in this project.
Every contribution is welcomed!
