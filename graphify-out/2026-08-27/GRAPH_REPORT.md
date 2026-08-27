# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 20 files · ~12,163 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 276 nodes · 593 edges · 34 communities (11 shown, 23 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 60 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f89ed9fc`
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
- NewChain
- sync.RWMutex
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- E
- Chain
- node
- E
- node
- E
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- Method
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
7. `assertFound()` - 11 edges
8. `setup()` - 11 edges
9. `mergeNodes()` - 10 edges
10. `RouterContext` - 10 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `TestIsParam()` --calls--> `isParam()`  [INFERRED]
  routing/routing_test.go → routing/routing.go
- `TestWithParams()` --calls--> `Param()`  [INFERRED]
  router/context/context_test.go → router/context/context.go
- `TestParam()` --calls--> `NewContext()`  [INFERRED]
  router/context/context_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (34 total, 23 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.09
Nodes (43): CreateTree(), assertFound(), assertPanicsWith(), Method, TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath() (+35 more)

### Community 2 - "testing.T"
Cohesion: 0.07
Nodes (63): net/http/httptest.Server, testing.T, Param(), TestParam(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_AnswersOptionsAutomatically() (+55 more)

### Community 3 - "FromRequest"
Cohesion: 0.15
Nodes (19): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, FromParams(), FromRequest() (+11 more)

### Community 6 - "Router"
Cohesion: 0.12
Nodes (14): net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer, sync.Mutex, Router, advertisedMethods(), defaultMethodNotAllowed() (+6 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 17 - "E"
Cohesion: 0.21
Nodes (14): Match, Param, classify(), E, Method, isParam(), nameParams(), splitSegments() (+6 more)

### Community 18 - "Chain"
Cohesion: 0.20
Nodes (8): Chain, Middleware, Router, net/http.Handler, NewPrefixRouter(), NewRouter(), NewPrefixRouter(), NewRouter()

### Community 20 - "node"
Cohesion: 0.23
Nodes (14): E, Method, Param, endpoint, lookup, node, cloneEndpoint(), cloneNode() (+6 more)

## Knowledge Gaps
- **5 isolated node(s):** `Install`, `Middleware`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `routing_test.go` to `E`, `Chain`?**
  _High betweenness centrality (0.296) - this node is a cross-community bridge._
- **Why does `NewPrefixRouter()` connect `Chain` to `routing_test.go`, `Router`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Why does `Tree` connect `E` to `routing_test.go`, `node`?**
  _High betweenness centrality (0.132) - this node is a cross-community bridge._
- **Are the 36 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 36 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Middleware`, `Instalação` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `routing_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09090909090909091 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06538461538461539 - nodes in this community are weakly interconnected._