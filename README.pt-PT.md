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
        
	    _const "github.com/joaolaureano/go-router/const"
        "github.com/joaolaureano/go-router/router"
        "github.com/joaolaureano/go-router/router/context"
    )
    
    func main() {
    r := router.NewRouter()
    
        r.Register(_const.GET, "/ping", func(writer http.ResponseWriter, request *http.Request) {
            writer.Write([]byte("pong"))
        })
        r.Group("/{id}", func(r router.Router) {
            r.Use(func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
                    requestCtx := request.Context().Value(context.RouterContextKey)
                    ctxValue, _ := requestCtx.(*context.RouterContext)
                    next.ServeHTTP(writer, request)
                    message := fmt.Sprintf("Group Middleware\n Found Path value: " + string(ctxValue.Value("id")))
                    writer.Write([]byte(message))
                })
            })
            r.Register(_const.GET, "/pong", func(writer http.ResponseWriter, request *http.Request) {
                writer.Write([]byte("ping"))
            })
        })
        http.ListenAndServe(":3333", r)
    }
```
Mais exemplos no diretório ```.example/```

## Interface
- `Register(httpMethod _const.HTTPMethods, path string, method http.HandlerFunc)`: Registra um método HTTP para um determinado caminho.
- `Use(middleware func(http.Handler) http.Handler)`: Utiliza um middleware para manipular as requisições HTTP.
- `NotFound(notFoundFn http.HandlerFunc)`: Define um handler para requisições em rotas não encontradas.
- `Group(prefix string, fn func(r router.Router)) router.Router`: Agrupa rotas com um determinado prefixo.
- `With(middleware ...func(http.Handler) http.Handler) *router.Router`: Utiliza middleware para um conjunto específico de rotas.

## Créditos

Este projeto foi inspirado e influenciado pelo **[go-Chi](https://github.com/go-chi/chi)**.

## Contribuindo

Sinta-se à vontade para abrir problemas ou enviar pull requests para contribuir com melhorias neste projeto. 
Toda contribuição é bem-vinda!