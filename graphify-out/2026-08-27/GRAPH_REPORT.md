# Graph Report - go-router  (2026-08-26)

## Corpus Check
- 14 files · ~4,522 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 120 nodes · 274 edges · 10 communities (8 shown, 2 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 46 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f3abbd8d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- Tree
- NewRouter
- RouterContext
- Router
- net/http.Handler
- chain_test.go
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 24 edges
2. `NewRouter()` - 20 edges
3. `Tree` - 17 edges
4. `Router` - 13 edges
5. `setup()` - 11 edges
6. `Node` - 9 edges
7. `RouterContext` - 8 edges
8. `Chain` - 8 edges
9. `isParam()` - 6 edges
10. `HTTPMethods` - 6 edges

## Surprising Connections (you probably didn't know these)
- `setPathVariableValues()` --references--> `RouterContext`  [EXTRACTED]
  tree/tree.go → router/context/context.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go

## Import Cycles
- None detected.

## Communities (10 total, 2 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.20
Nodes (26): testing.T, CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath(), TestFindRoute_InexistentPathOnlyRoot() (+18 more)

### Community 1 - "Tree"
Cohesion: 0.26
Nodes (9): HTTPMethods, Method, Node, RouterTree, Tree, isParam(), namedPath(), setPathVariableValues() (+1 more)

### Community 2 - "NewRouter"
Cohesion: 0.22
Nodes (19): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_Group(), TestRouter_MiddlewareOrder(), TestRouter_NestedGroupsComposePrefixes(), TestRouter_NotFound() (+11 more)

### Community 3 - "RouterContext"
Cohesion: 0.19
Nodes (9): net/http.Request, net/http.ResponseWriter, RouterContext, NewContext(), setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey(), TestContext_ValidKey() (+1 more)

### Community 4 - "Router"
Cohesion: 0.23
Nodes (8): Middleware, _const.HTTPMethods, net/http.HandlerFunc, sync.RWMutex, Router, NewPrefixRouter(), TestNewRouterWithPrefix(), TestRouter_WithPreservesPrefix()

### Community 5 - "net/http.Handler"
Cohesion: 0.43
Nodes (4): Chain, Router, net/http.Handler, NewRouter()

### Community 6 - "chain_test.go"
Cohesion: 0.28
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

## Knowledge Gaps
- **1 isolated node(s):** `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `Tree`, `NewRouter`, `Router`?**
  _High betweenness centrality (0.186) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `testing.T`?**
  _High betweenness centrality (0.152) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `Router`, `net/http.Handler`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestFindRoute()`) actually correct?**
  _`CreateTree()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 16 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_Group()`) actually correct?**
  _`NewRouter()` has 16 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/joaolaureano/go-router` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._