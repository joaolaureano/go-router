# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~7,211 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 189 nodes · 448 edges · 19 communities (10 shown, 9 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 84 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b4d47618`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- node
- NewRouter
- FromRequest
- _const.HTTPMethods
- methodFilter
- net/http.Handler
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- Go-Router
- Node
- Tree
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- newNode
- Router
- NewChain

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 36 edges
2. `CreateTree()` - 29 edges
3. `Router` - 23 edges
4. `node` - 18 edges
5. `Tree` - 16 edges
6. `Method` - 13 edges
7. `setup()` - 11 edges
8. `FromRequest()` - 11 edges
9. `assertFound()` - 9 edges
10. `Chain` - 9 edges

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

## Communities (19 total, 9 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.19
Nodes (30): testing.T, CreateTree(), assertFound(), TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+22 more)

### Community 1 - "node"
Cohesion: 0.28
Nodes (7): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "NewRouter"
Cohesion: 0.12
Nodes (34): net/http/httptest.Server, NewPrefixRouter(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_Group(), TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute() (+26 more)

### Community 3 - "FromRequest"
Cohesion: 0.13
Nodes (22): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param() (+14 more)

### Community 6 - "net/http.Handler"
Cohesion: 0.22
Nodes (7): Chain, Middleware, Router, github.com/joaolaureano/go-router/router.Router, net/http.Handler, NewPrefixRouter(), NewRouter()

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.29
Nodes (10): Match, Param, Status, Tree, isParam(), splitSegments(), TestValidatePath_InvalidPaths(), TestValidatePath_ValidPaths() (+2 more)

### Community 16 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

### Community 17 - "Router"
Cohesion: 0.30
Nodes (5): net/http.HandlerFunc, sync.RWMutex, Method, Router, routeTree

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

## Knowledge Gaps
- **5 isolated node(s):** `Install`, `Instalação`, `routeTree`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Router` connect `Router` to `NewRouter`, `NewChain`, `FromRequest`, `net/http.Handler`?**
  _High betweenness centrality (0.165) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `Router`, `FromRequest`, `net/http.Handler`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Why does `CreateTree()` connect `testing.T` to `newNode`, `NewRouter`, `Tree`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Are the 33 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `routeTree` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewRouter` be split into smaller, more focused modules?**
  _Cohesion score 0.11596638655462185 - nodes in this community are weakly interconnected._