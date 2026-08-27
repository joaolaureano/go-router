# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 19 files · ~11,855 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 277 nodes · 595 edges · 33 communities (10 shown, 23 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 60 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `58840a5f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- routing_test.go
- tree.Method
- testing.T
- FromRequest
- _const.HTTPMethods
- methodFilter
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- Go-Router
- Chain
- sync.RWMutex
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- E
- Method
- node
- E
- node
- E
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- segmentKind
- E
- tree.Param
- Method

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 40 edges
2. `Router` - 23 edges
3. `node` - 20 edges
4. `Param()` - 14 edges
5. `FromRequest()` - 12 edges
6. `state` - 12 edges
7. `setup()` - 11 edges
8. `assertFound()` - 11 edges
9. `node[E]` - 10 edges
10. `mergeNodes()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `TestWithParams()` --calls--> `Param()`  [INFERRED]
  router/context/context_test.go → router/context/context.go
- `TestParam()` --calls--> `NewContext()`  [INFERRED]
  router/context/context_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (33 total, 23 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.09
Nodes (42): CreateTree(), assertFound(), assertPanicsWith(), Method, TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+34 more)

### Community 2 - "testing.T"
Cohesion: 0.07
Nodes (63): net/http/httptest.Server, testing.T, Param(), TestParam(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_AnswersOptionsAutomatically() (+55 more)

### Community 3 - "FromRequest"
Cohesion: 0.15
Nodes (19): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, FromParams(), FromRequest() (+11 more)

### Community 6 - "Router"
Cohesion: 0.10
Nodes (19): Router, net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer, sync.Mutex, NewPrefixRouter(), NewRouter() (+11 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.15
Nodes (10): Chain, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares() (+2 more)

### Community 17 - "E"
Cohesion: 0.20
Nodes (15): Match, Param, classify(), E, Method, isParam(), nameParams(), splitSegments() (+7 more)

### Community 20 - "node"
Cohesion: 0.22
Nodes (14): E, Param, endpoint, lookup, node, cloneEndpoint(), cloneNode(), node[E] (+6 more)

## Knowledge Gaps
- **5 isolated node(s):** `Middleware`, `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `routing_test.go` to `E`, `Router`?**
  _High betweenness centrality (0.299) - this node is a cross-community bridge._
- **Why does `NewPrefixRouter()` connect `Router` to `routing_test.go`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Why does `Tree` connect `E` to `routing_test.go`, `node`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Are the 36 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 36 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Middleware`, `Install`, `Instalação` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `routing_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09413067552602436 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06538461538461539 - nodes in this community are weakly interconnected._