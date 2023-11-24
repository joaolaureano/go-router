# Go-Router

Go-Router é um roteador em Golang simples para requisições HTTP. 
O projeto foi desenvolvido com o fim de aprendizado e experimentação da linguagem Go.

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