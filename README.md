[![en](https://img.shields.io/badge/lang-en-blue.svg)](https://github.com/joaolaureano/go-router/blob/main/README.md)
[![pt-br](https://img.shields.io/badge/lang-pt--br-green.svg)](https://github.com/joaolaureano/go-router/blob/main/README.pt-BR.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/joaolaureano/go-router.svg)](https://pkg.go.dev/github.com/joaolaureano/go-router)
[![Go Report Card](https://goreportcard.com/badge/github.com/joaolaureano/go-router)](https://goreportcard.com/report/github.com/joaolaureano/go-router)

# go-router

A small, dependency-free HTTP router for Go, built on a compressed (patricia) trie with lock-free lookups. Started as a learning project inspired by [chi](https://github.com/go-chi/chi); it has since grown a fuzz-tested routing tree and a [load-tested](loadtest/README.md) concurrency story.

## Features

- **Patricia trie matching** — static routes and `{param}` captures only, like a standard router. A shared literal prefix, even across a `/`, collapses into one edge instead of one node per segment, splitting only at the byte where two routes first diverge.
- **Groups & mounting** — `Group` for shared prefixes and middleware, `Mount` for grafting one router's routes under another's prefix.
- **Middleware chains** — `Use` for a router or group, `With` for a one-off set of middleware on specific routes.
- **Correct method fallbacks** — `HEAD` falls back to `GET`, `OPTIONS` is answered automatically with a computed `Allow` header, unless you register your own.
- **Customizable fallbacks** — `NotFound` and `MethodNotAllowed` handlers.
- **Concurrency-safe by design** — the routing tree is swapped atomically, so registering routes while serving traffic never exposes a half-built tree. Verified under load, not just asserted (see [`loadtest/`](loadtest/README.md)).
- **Zero runtime dependencies** — only the standard library; `testify` is a test-only dependency.

## Install

```sh
go get github.com/joaolaureano/go-router@latest
```

Requires Go 1.26 or later (see [go.mod](go.mod)).

## Quickstart

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

    r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

    r.Group("/{id}", func(r *router.Router) {
        r.Use(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                next.ServeHTTP(w, req)
                message := fmt.Sprintf("Group middleware, found id: %s", context.Param(req, "id"))
                w.Write([]byte(message))
            })
        })
        r.Get("/pong", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("ping"))
        })
    })

    http.ListenAndServe(":3333", r)
}
```

More runnable examples in [`.example/`](.example/).

## Routing patterns

### Path parameters

A `{name}` segment captures whatever the request has in that position, read back with `context.Param`:

```go
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := context.Param(r, "id")
    w.Write([]byte("user " + id))
})
```

### Method handling

`HEAD` falls back to the `GET` route unless one is registered for it, and `OPTIONS` on a known path answers `204` with an `Allow` header unless a route claims it. Registering either explicitly always wins. `Allow` advertises the registered methods plus the two the router answers on their behalf.

## Interface

| Method | Purpose |
|---|---|
| `Register(method router.Method, path string, handler http.HandlerFunc)` | Registers an HTTP method for a specific path. |
| `Get/Head/Post/Put/Patch/Delete/Options(path string, handler http.HandlerFunc)` | Shortcuts for `Register` with the verb spelled into the name. |
| `Use(middleware func(http.Handler) http.Handler)` | Adds middleware applied to every route on this router or group. |
| `With(middleware ...func(http.Handler) http.Handler) *router.Router` | Applies middleware to a specific set of routes without affecting the rest. |
| `Group(prefix string, fn func(r *router.Router)) *router.Router` | Groups routes under a shared prefix, with its own middleware. |
| `Mount(prefix string, other *router.Router)` | Grafts another router's routes under a prefix, handlers and middleware intact. |
| `NotFound(handler http.HandlerFunc)` | Sets a handler for requests on non-existent routes. |
| `MethodNotAllowed(handler http.HandlerFunc)` | Sets a handler for requests to an existing path under an unregistered method. |

## Testing & performance

The routing tree carries three layers of verification beyond ordinary unit tests:

- **Fuzz tests** (`routing/fuzz_test.go`, `router/fuzz_test.go`) exercise route registration and lookup against arbitrary inputs.
- **Benchmarks** (`router/bench_test.go`) cover static, parameter, escaped, and miss lookups, including parallel and large-fanout cases.
- **Load tests** (`loadtest/`), driven by [vegeta](https://github.com/tsenart/vegeta), verify correctness *under concurrent load* — including registering or mounting routes while traffic is in flight. See [`loadtest/README.md`](loadtest/README.md) for the full scenario list and how to run them.

```sh
go test ./...                     # unit + fuzz-as-unit tests
go test -bench=. -run=^$ ./router # benchmarks
cd loadtest && go test ./...      # load tests (separate module)
```

## Credits

This project was inspired and influenced by **[chi](https://github.com/go-chi/chi)**.

## Contributing

Feel free to open issues or send pull requests to contribute to improvements in this project. Every contribution is welcome!
