# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,772 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 164 nodes · 389 edges · 13 communities (9 shown, 4 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 77 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `dfbf06a0`
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
- net/http.Handler
- github.com/joaolaureano/go-router/router/context.RouterContext

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `NewRouter()` - 28 edges
3. `Node` - 20 edges
4. `Tree` - 16 edges
5. `Router` - 14 edges
6. `HTTPMethods` - 11 edges
7. `setup()` - 11 edges
8. `NewContext()` - 11 edges
9. `Chain` - 9 edges
10. `newNode()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `TestFindRoute_FallsBackToParameterAfterStaticBranchFails()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go
- `TestTree_MergeCopiesRoutesWithoutRebuildingPaths()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (13 total, 4 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.18
Nodes (29): testing.T, CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_FallsBackToParameterAfterStaticBranchFails(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath() (+21 more)

### Community 1 - "Node"
Cohesion: 0.22
Nodes (11): HTTPMethods, Method, methodFilter, Node, cloneMethod(), cloneNode(), filterByAnyMethod(), filterByMethod() (+3 more)

### Community 2 - "NewRouter"
Cohesion: 0.15
Nodes (28): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group(), TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute(), TestRouter_GroupNotFoundReachesServingRouter() (+20 more)

### Community 3 - "NewContext"
Cohesion: 0.17
Nodes (14): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+6 more)

### Community 5 - "Tree"
Cohesion: 0.24
Nodes (10): methodFilter, Node, RouterTree, Tree, isParam(), setPathVariableValues(), splitSegments(), TestValidatePath_InvalidPaths() (+2 more)

### Community 6 - "Router"
Cohesion: 0.12
Nodes (15): Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares() (+7 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

### Community 11 - "net/http.Handler"
Cohesion: 0.36
Nodes (4): Chain, Router, net/http.Handler, NewRouter()

## Knowledge Gaps
- **2 isolated node(s):** `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `NewRouter`, `newNode`, `Tree`, `Router`?**
  _High betweenness centrality (0.152) - this node is a cross-community bridge._
- **Why does `Node` connect `Node` to `newNode`?**
  _High betweenness centrality (0.112) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `net/http.Handler`, `Router`?**
  _High betweenness centrality (0.107) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewRouter` be split into smaller, more focused modules?**
  _Cohesion score 0.1477832512315271 - nodes in this community are weakly interconnected._