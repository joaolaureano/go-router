# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,055 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 138 nodes · 320 edges · 11 communities (9 shown, 2 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 60 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e02d1b34`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- Node
- NewRouter
- RouterContext
- net/http.Handler
- Router
- chain_test.go
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `NewRouter()` - 21 edges
3. `Node` - 19 edges
4. `Tree` - 15 edges
5. `Router` - 13 edges
6. `setup()` - 11 edges
7. `newNode()` - 8 edges
8. `NewContext()` - 8 edges
9. `RouterContext` - 8 edges
10. `Chain` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `setPathVariableValues()` --references--> `RouterContext`  [EXTRACTED]
  tree/tree.go → router/context/context.go
- `TestFindRoute_FallsBackToParameterAfterStaticBranchFails()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go
- `TestTree_MergeCopiesRoutesWithoutRebuildingPaths()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (11 total, 2 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.17
Nodes (31): testing.T, CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_FallsBackToParameterAfterStaticBranchFails(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath() (+23 more)

### Community 1 - "Node"
Cohesion: 0.17
Nodes (13): HTTPMethods, Method, Node, cloneMethod(), cloneNode(), matchPath(), mergeNodes(), staticChild() (+5 more)

### Community 2 - "NewRouter"
Cohesion: 0.21
Nodes (20): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_Group(), TestRouter_MiddlewareOrder(), TestRouter_NestedGroupsComposePrefixes(), TestRouter_NotFound() (+12 more)

### Community 3 - "RouterContext"
Cohesion: 0.19
Nodes (9): net/http.Request, net/http.ResponseWriter, RouterContext, NewContext(), setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey(), TestContext_ValidKey() (+1 more)

### Community 4 - "net/http.Handler"
Cohesion: 0.43
Nodes (4): Chain, Router, net/http.Handler, NewRouter()

### Community 5 - "Router"
Cohesion: 0.23
Nodes (8): Middleware, net/http.HandlerFunc, sync.RWMutex, Router, _const.HTTPMethods, NewPrefixRouter(), TestNewRouterWithPrefix(), TestRouter_WithPreservesPrefix()

### Community 6 - "chain_test.go"
Cohesion: 0.28
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

## Knowledge Gaps
- **1 isolated node(s):** `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `Node`, `NewRouter`, `newNode`, `Router`?**
  _High betweenness centrality (0.191) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `net/http.Handler`, `Router`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Why does `Tree` connect `Node` to `testing.T`?**
  _High betweenness centrality (0.119) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_Group()`) actually correct?**
  _`NewRouter()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/joaolaureano/go-router` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._