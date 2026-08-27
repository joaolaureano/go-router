# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 19 files · ~11,485 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 273 nodes · 639 edges · 34 communities (10 shown, 24 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 112 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d2473922`
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
- E
- node
- E
- node
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- Method
- segmentKind
- Param
- E
- tree.Param
- Method

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 55 edges
2. `CreateTree()` - 38 edges
3. `Router` - 24 edges
4. `node` - 19 edges
5. `Param()` - 14 edges
6. `FromRequest()` - 12 edges
7. `state` - 12 edges
8. `assertFound()` - 11 edges
9. `setup()` - 11 edges
10. `node[E]` - 10 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `TestIsParam()` --calls--> `isParam()`  [INFERRED]
  routing/routing_test.go → routing/routing.go

## Import Cycles
- None detected.

## Communities (34 total, 24 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.12
Nodes (45): testing.T, TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint(), CreateTree(), assertFound(), assertPanicsWith() (+37 more)

### Community 2 - "NewRouter"
Cohesion: 0.08
Nodes (53): net/http/httptest.Server, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_AnswersOptionsAutomatically(), TestRouter_CatchAllRoute(), TestRouter_CatchAllUnderGroup() (+45 more)

### Community 3 - "FromRequest"
Cohesion: 0.12
Nodes (23): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, FromParams(), FromRequest() (+15 more)

### Community 6 - "Router"
Cohesion: 0.11
Nodes (17): net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer, sync.Mutex, Method, Router, advertisedMethods() (+9 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.12
Nodes (13): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+5 more)

### Community 17 - "E"
Cohesion: 0.21
Nodes (14): Match, Param, classify(), E, Method, isParam(), nameParams(), splitSegments() (+6 more)

### Community 20 - "node"
Cohesion: 0.23
Nodes (12): endpoint, lookup, node, cloneEndpoint(), cloneNode(), node[E], E, Method (+4 more)

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `E`, `Router`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `Tree` connect `E` to `testing.T`, `node`?**
  _High betweenness centrality (0.144) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewRouter`, `Chain`?**
  _High betweenness centrality (0.134) - this node is a cross-community bridge._
- **Are the 52 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_AnswersOptionsAutomatically()`) actually correct?**
  _`NewRouter()` has 52 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 34 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.11748381128584644 - nodes in this community are weakly interconnected._