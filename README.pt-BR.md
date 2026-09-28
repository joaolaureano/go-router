[![en](https://img.shields.io/badge/lang-en-blue.svg)](https://github.com/joaolaureano/go-router/blob/main/README.md)
[![pt-br](https://img.shields.io/badge/lang-pt--br-green.svg)](https://github.com/joaolaureano/go-router/blob/main/README.pt-BR.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/joaolaureano/go-router.svg)](https://pkg.go.dev/github.com/joaolaureano/go-router)
[![Go Report Card](https://goreportcard.com/badge/github.com/joaolaureano/go-router)](https://goreportcard.com/report/github.com/joaolaureano/go-router)

# go-router

Um roteador HTTP pequeno e sem dependências para Go, construído sobre uma patricia trie (radix tree comprimida) com buscas lock-free. Começou como um projeto de aprendizado inspirado no [chi](https://github.com/go-chi/chi); desde então ganhou uma árvore de roteamento testada com fuzzing e uma história de concorrência validada por [teste de carga](loadtest/README.md).

## Funcionalidades

- **Matching por patricia trie** — só rotas estáticas e captura via `{param}`, como um router padrão. Um prefixo literal compartilhado, mesmo atravessando uma `/`, vira uma única aresta em vez de um nó por segmento, dividindo-se só no byte em que duas rotas divergem.
- **Grupos e montagem** — `Group` para prefixos e middlewares compartilhados, `Mount` para enxertar as rotas de um router sob o prefixo de outro.
- **Cadeias de middleware** — `Use` para um router ou grupo, `With` para aplicar middleware a um conjunto específico de rotas.
- **Fallback correto de métodos** — `HEAD` recorre ao `GET`, `OPTIONS` é respondido automaticamente com um `Allow` calculado, a não ser que você registre o seu próprio.
- **Fallbacks customizáveis** — handlers de `NotFound` e `MethodNotAllowed`.
- **Seguro para concorrência por design** — a árvore de rotas é trocada atomicamente, então registrar rotas enquanto o servidor está no ar nunca expõe uma árvore pela metade. Verificado sob carga, não apenas assumido (veja [`loadtest/`](loadtest/README.md)).
- **Zero dependências em produção** — só a standard library; `testify` é dependência exclusiva de teste.

## Instalação

```sh
go get github.com/joaolaureano/go-router@latest
```

Requer Go 1.26 ou superior (veja [go.mod](go.mod)).

## Início rápido

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
                message := fmt.Sprintf("Middleware do grupo, id encontrado: %s", context.Param(req, "id"))
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

Mais exemplos executáveis em [`.example/`](.example/).

## Padrões de rota

### Parâmetros de caminho

Um segmento `{name}` captura o que estiver naquela posição da requisição, lido de volta com `context.Param`:

```go
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := context.Param(r, "id")
    w.Write([]byte("user " + id))
})
```

### Tratamento de métodos

`HEAD` recorre à rota `GET` a não ser que exista uma registrada para ele, e `OPTIONS` num caminho conhecido responde `204` com cabeçalho `Allow` a não ser que uma rota o reclame. Registrar qualquer um deles explicitamente tem sempre precedência. O `Allow` anuncia os métodos registrados mais os dois que o router responde por conta deles.

## Interface

| Método | Propósito |
|---|---|
| `Register(method router.Method, path string, handler http.HandlerFunc)` | Registra um método HTTP para um determinado caminho. |
| `Get/Head/Post/Put/Patch/Delete/Options(path string, handler http.HandlerFunc)` | Atalhos para `Register` com o verbo no nome. |
| `Use(middleware func(http.Handler) http.Handler)` | Adiciona middleware aplicado a toda rota deste router ou grupo. |
| `With(middleware ...func(http.Handler) http.Handler) *router.Router` | Aplica middleware a um conjunto específico de rotas, sem afetar as demais. |
| `Group(prefix string, fn func(r *router.Router)) *router.Router` | Agrupa rotas sob um prefixo compartilhado, com middleware próprio. |
| `Mount(prefix string, other *router.Router)` | Enxerta as rotas de outro router sob um prefixo, mantendo handlers e middleware. |
| `NotFound(handler http.HandlerFunc)` | Define um handler para requisições em rotas não encontradas. |
| `MethodNotAllowed(handler http.HandlerFunc)` | Define um handler para requisições a um caminho existente sob um método não registrado. |

## Testes e performance

A árvore de roteamento carrega três camadas de verificação além dos testes unitários comuns:

- **Testes de fuzzing** (`routing/fuzz_test.go`, `router/fuzz_test.go`) exercitam registro e busca de rotas contra entradas arbitrárias.
- **Benchmarks** (`router/bench_test.go`) cobrem buscas estáticas, com parâmetro, com escape e sem match, incluindo casos paralelos e de alto fan-out.
- **Testes de carga** (`loadtest/`), rodados com [vegeta](https://github.com/tsenart/vegeta), verificam a corretude *sob carga concorrente* — incluindo registrar ou montar rotas com tráfego em andamento. Veja [`loadtest/README.md`](loadtest/README.md) para a lista completa de cenários e como rodá-los.

```sh
go test ./...                     # testes unitários + fuzzing como teste unitário
go test -bench=. -run=^$ ./router # benchmarks
cd loadtest && go test ./...      # testes de carga (módulo separado)
```

## Créditos

Este projeto foi inspirado e influenciado pelo **[chi](https://github.com/go-chi/chi)**.

## Contribuindo

Sinta-se à vontade para abrir issues ou enviar pull requests para contribuir com melhorias neste projeto. Toda contribuição é bem-vinda!
