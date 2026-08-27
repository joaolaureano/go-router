# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 19 files · ~11,896 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 277 nodes · 650 edges · 34 communities (11 shown, 23 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 114 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f71f7309`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- tree.Method
- NewRouter
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
- NewChain
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
1. `NewRouter()` - 55 edges
2. `CreateTree()` - 40 edges
3. `Router` - 24 edges
4. `node` - 20 edges
5. `Param()` - 14 edges
6. `FromRequest()` - 12 edges
7. `state` - 12 edges
8. `assertFound()` - 11 edges
9. `setup()` - 11 edges
10. `node[E]` - 10 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `TestIsParam()` --calls--> `isParam()`  [INFERRED]
  routing/routing_test.go → routing/routing.go

## Import Cycles
- None detected.

## Communities (34 total, 23 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.11
Nodes (48): testing.T, TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint(), CreateTree(), assertFound(), assertPanicsWith() (+40 more)

### Community 2 - "NewRouter"
Cohesion: 0.08
Nodes (53): net/http/httptest.Server, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_AnswersOptionsAutomatically(), TestRouter_CatchAllRoute(), TestRouter_CatchAllUnderGroup() (+45 more)

### Community 3 - "FromRequest"
Cohesion: 0.12
Nodes (23): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, FromParams(), FromRequest() (+15 more)

### Community 6 - "Router"
Cohesion: 0.12
Nodes (14): net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer, sync.Mutex, Method, Router, advertisedMethods() (+6 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.17
Nodes (9): Chain, Middleware, Router, net/http.Handler, NewPrefixRouter(), NewRouter(), NewPrefixRouter(), TestNewRouterWithPrefix() (+1 more)

### Community 17 - "E"
Cohesion: 0.21
Nodes (14): Match, Param, classify(), E, Method, isParam(), nameParams(), splitSegments() (+6 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 20 - "node"
Cohesion: 0.22
Nodes (14): E, Param, endpoint, lookup, node, cloneEndpoint(), cloneNode(), node[E] (+6 more)

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `E`, `Chain`?**
  _High betweenness centrality (0.250) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewRouter`, `NewChain`, `Chain`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Why does `Tree` connect `E` to `testing.T`, `node`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Are the 52 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_AnswersOptionsAutomatically()`) actually correct?**
  _`NewRouter()` has 52 INFERRED edges - model-reasoned connections that need verification._
- **Are the 36 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 36 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.10938775510204081 - nodes in this community are weakly interconnected._