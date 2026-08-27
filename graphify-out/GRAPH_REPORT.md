# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 18 files · ~9,096 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 223 nodes · 536 edges · 27 communities (11 shown, 16 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 93 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `56ca4b3c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- node
- NewRouter
- FromRequest
- _const.HTTPMethods
- methodFilter
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- Go-Router
- net/http.Handler
- Tree
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- Method
- NewChain
- node
- Method
- segmentKind
- Method
- Method
- Method

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 43 edges
2. `CreateTree()` - 36 edges
3. `Router` - 24 edges
4. `node` - 19 edges
5. `Tree` - 16 edges
6. `Method` - 14 edges
7. `setup()` - 11 edges
8. `FromRequest()` - 11 edges
9. `assertFound()` - 10 edges
10. `Param()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `TestRouter_BacktracksWhenStaticMatchLacksMethod()` --calls--> `NewRouter()`  [INFERRED]
  router/router_test.go → router/router.go

## Import Cycles
- None detected.

## Communities (27 total, 16 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.12
Nodes (43): testing.T, TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint(), CreateTree(), assertFound(), assertPanicsWith() (+35 more)

### Community 1 - "node"
Cohesion: 0.28
Nodes (12): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), node[E], E (+4 more)

### Community 2 - "NewRouter"
Cohesion: 0.10
Nodes (41): net/http/httptest.Server, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_CatchAllRoute(), TestRouter_CatchAllUnderGroup(), TestRouter_Group() (+33 more)

### Community 3 - "FromRequest"
Cohesion: 0.17
Nodes (16): contextKey, param, RouterContext, FromRequest(), NewContext(), setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey() (+8 more)

### Community 6 - "Router"
Cohesion: 0.21
Nodes (7): net/http.HandlerFunc, net/http.Request, net/http.ResponseWriter, sync.RWMutex, Router, defaultMethodNotAllowed(), routes

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "net/http.Handler"
Cohesion: 0.17
Nodes (9): Chain, Middleware, Router, net/http.Handler, NewPrefixRouter(), NewRouter(), NewPrefixRouter(), TestNewRouterWithPrefix() (+1 more)

### Community 12 - "Tree"
Cohesion: 0.25
Nodes (13): Match, Param, Status, Tree, classify(), Tree[E], E, isParam() (+5 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `assertFound()` connect `testing.T` to `node`, `net/http.Handler`, `Tree`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Why does `Method` connect `node` to `testing.T`, `Tree`, `Router`?**
  _High betweenness centrality (0.124) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewRouter`, `NewChain`, `net/http.Handler`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Are the 40 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 40 INFERRED edges - model-reasoned connections that need verification._
- **Are the 33 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.12323232323232323 - nodes in this community are weakly interconnected._