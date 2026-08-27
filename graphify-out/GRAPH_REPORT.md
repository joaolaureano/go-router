# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,168 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 145 nodes · 339 edges · 12 communities (9 shown, 3 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 65 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e7dc6d33`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- Node
- NewRouter
- NewContext
- _const.HTTPMethods
- Router
- chain_test.go
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode
- net/http.Handler

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `NewRouter()` - 21 edges
3. `Node` - 19 edges
4. `Tree` - 15 edges
5. `Router` - 14 edges
6. `NewContext()` - 11 edges
7. `setup()` - 11 edges
8. `HTTPMethods` - 9 edges
9. `RouterContext` - 8 edges
10. `newNode()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `TestFindRoute_FallsBackToParameterAfterStaticBranchFails()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go
- `TestTree_MergeCopiesRoutesWithoutRebuildingPaths()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go
- `TestTree_MergeKeepsExistingRoute()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go
- `TestTree_MergeWithItselfIsNoOp()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go

## Import Cycles
- None detected.

## Communities (12 total, 3 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.17
Nodes (31): testing.T, CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_FallsBackToParameterAfterStaticBranchFails(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath() (+23 more)

### Community 1 - "Node"
Cohesion: 0.17
Nodes (13): HTTPMethods, Method, Node, cloneMethod(), cloneNode(), matchPath(), mergeNodes(), staticChild() (+5 more)

### Community 2 - "NewRouter"
Cohesion: 0.21
Nodes (20): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_Group(), TestRouter_MiddlewareOrder(), TestRouter_NestedGroupsComposePrefixes(), TestRouter_NotFound() (+12 more)

### Community 3 - "NewContext"
Cohesion: 0.17
Nodes (14): contextKey, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+6 more)

### Community 5 - "Router"
Cohesion: 0.23
Nodes (8): Middleware, github.com/joaolaureano/go-router/tree.RouterTree, net/http.HandlerFunc, sync.RWMutex, Router, NewPrefixRouter(), TestNewRouterWithPrefix(), TestRouter_WithPreservesPrefix()

### Community 6 - "chain_test.go"
Cohesion: 0.28
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

### Community 11 - "net/http.Handler"
Cohesion: 0.43
Nodes (4): Chain, Router, net/http.Handler, NewRouter()

## Knowledge Gaps
- **2 isolated node(s):** `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `Node`, `NewRouter`, `newNode`, `Router`?**
  _High betweenness centrality (0.169) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `net/http.Handler`, `Router`?**
  _High betweenness centrality (0.106) - this node is a cross-community bridge._
- **Why does `Tree` connect `Node` to `testing.T`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_Group()`) actually correct?**
  _`NewRouter()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._