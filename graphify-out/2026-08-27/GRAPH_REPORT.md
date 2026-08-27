# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,391 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 154 nodes · 367 edges · 12 communities (8 shown, 4 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 70 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e2856ac4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- Node
- NewRouter
- NewContext
- _const.HTTPMethods
- Tree
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode
- github.com/joaolaureano/go-router/tree.RouterTree

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `Node` - 24 edges
3. `NewRouter()` - 23 edges
4. `Tree` - 15 edges
5. `Router` - 14 edges
6. `setup()` - 11 edges
7. `HTTPMethods` - 11 edges
8. `NewContext()` - 11 edges
9. `newNode()` - 8 edges
10. `Chain` - 8 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `TestFindRoute_FallsBackToParameterAfterStaticBranchFails()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (12 total, 4 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.17
Nodes (31): testing.T, CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_FallsBackToParameterAfterStaticBranchFails(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath() (+23 more)

### Community 1 - "Node"
Cohesion: 0.19
Nodes (13): HTTPMethods, github.com/joaolaureano/go-router/router/context.RouterContext, Method, methodFilter, Node, cloneMethod(), cloneNode(), filterByAnyMethod() (+5 more)

### Community 2 - "NewRouter"
Cohesion: 0.15
Nodes (26): net/http/httptest.Server, NewPrefixRouter(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+18 more)

### Community 3 - "NewContext"
Cohesion: 0.17
Nodes (14): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+6 more)

### Community 5 - "Tree"
Cohesion: 0.21
Nodes (9): Chain, Middleware, Router, net/http.Handler, NewRouter(), tree.RouterTree, Tree, isParam() (+1 more)

### Community 6 - "Router"
Cohesion: 0.17
Nodes (10): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares(), net/http.HandlerFunc (+2 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

## Knowledge Gaps
- **2 isolated node(s):** `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `NewRouter`, `newNode`, `Tree`?**
  _High betweenness centrality (0.159) - this node is a cross-community bridge._
- **Why does `Node` connect `Node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.127) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `testing.T`, `Node`?**
  _High betweenness centrality (0.106) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._