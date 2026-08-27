# Go-Router

Go-Router é um roteador em Golang simples para requisições HTTP. 
O projeto foi desenvolvido com o fim de aprendizado e experimentação da linguagem Go.

## Instalação
```go get github.com/joaolaureano/go-router@latest```


### Exemplo de código
**Tão simples como abaixo**
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
Mais exemplos no diretório ```.example/```

## Interface
- `Register(httpMethod router.Method, path string, method http.HandlerFunc)`: Registra um método HTTP para um determinado caminho.
- `Get/Head/Post/Put/Patch/Delete/Options(path string, handler http.HandlerFunc)`: Atalhos para `Register` com o verbo no nome.
- `Use(middleware func(http.Handler) http.Handler)`: Utiliza um middleware para manipular as requisições HTTP.
- `NotFound(notFoundFn http.HandlerFunc)`: Define um handler para requisições em rotas não encontradas.
- `MethodNotAllowed(fn http.HandlerFunc)`: Sets a handler for requests to an existing path under an unregistered method.
- `Group(prefix string, fn func(r *router.Router)) *router.Router`: Agrupa rotas com um determinado prefixo.
- `Mount(prefix string, other *router.Router)`: Enxerta as rotas de outro router sob um prefixo, mantendo handlers e middleware.
- `With(middleware ...func(http.Handler) http.Handler) *router.Router`: Utiliza middleware para um conjunto específico de rotas.

## Créditos

Este projeto foi inspirado e influenciado pelo **[go-Chi](https://github.com/go-chi/chi)**.

## Contribuindo

Sinta-se à vontade para abrir problemas ou enviar pull requests para contribuir com melhorias neste projeto. 
Toda contribuição é bem-vinda!