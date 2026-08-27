# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 20 files · ~12,352 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 276 nodes · 594 edges · 32 communities (10 shown, 22 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 60 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b5e9b057`
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
- routing.go
- E
- E
- E
- E
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- Method
- E
- tree.Param
- Method

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 40 edges
2. `Router` - 24 edges
3. `node` - 18 edges
4. `Param()` - 14 edges
5. `FromRequest()` - 12 edges
6. `state` - 12 edges
7. `assertFound()` - 11 edges
8. `setup()` - 11 edges
9. `mergeNodes()` - 10 edges
10. `RouterContext` - 10 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `advertisedMethods()` --references--> `Method`  [EXTRACTED]
  router/router.go → routing/method.go
- `TestWithParams()` --calls--> `Param()`  [INFERRED]
  router/context/context_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (32 total, 22 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.09
Nodes (42): CreateTree(), assertFound(), assertPanicsWith(), Method, TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+34 more)

### Community 2 - "testing.T"
Cohesion: 0.07
Nodes (63): net/http/httptest.Server, testing.T, Param(), TestParam(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_AnswersOptionsAutomatically() (+55 more)

### Community 3 - "FromRequest"
Cohesion: 0.12
Nodes (23): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, net/http.ResponseWriter, FromParams() (+15 more)

### Community 6 - "Router"
Cohesion: 0.11
Nodes (15): Router, net/http.HandlerFunc, sync/atomic.Bool, sync/atomic.Pointer, sync.Mutex, NewPrefixRouter(), NewRouter(), Router (+7 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.15
Nodes (10): Chain, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares() (+2 more)

### Community 17 - "routing.go"
Cohesion: 0.22
Nodes (13): Method, Match, Param, classify(), isParam(), nameParams(), splitSegments(), TestIsParam() (+5 more)

### Community 20 - "E"
Cohesion: 0.17
Nodes (17): E, node, Param, endpoint, lookup, node, cloneEndpoint(), cloneNode() (+9 more)

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `routing_test.go` to `routing.go`, `E`, `Router`?**
  _High betweenness centrality (0.309) - this node is a cross-community bridge._
- **Why does `NewPrefixRouter()` connect `Router` to `routing_test.go`?**
  _High betweenness centrality (0.140) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `FromRequest`, `Chain`?**
  _High betweenness centrality (0.126) - this node is a cross-community bridge._
- **Are the 36 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 36 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `routing_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09413067552602436 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06538461538461539 - nodes in this community are weakly interconnected._