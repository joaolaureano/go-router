# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 14 files · ~4,796 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 124 nodes · 284 edges · 10 communities (8 shown, 2 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 49 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fad06ebb`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
- Tree
- testing.T
- RouterContext
- Router
- net/http.Handler
- chain_test.go
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 26 edges
2. `NewRouter()` - 21 edges
3. `Tree` - 18 edges
4. `Router` - 13 edges
5. `setup()` - 11 edges
6. `Node` - 10 edges
7. `RouterContext` - 8 edges
8. `Chain` - 8 edges
9. `HTTPMethods` - 6 edges
10. `setup()` - 6 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `setPathVariableValues()` --references--> `RouterContext`  [EXTRACTED]
  tree/tree.go → router/context/context.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go

## Import Cycles
- None detected.

## Communities (10 total, 2 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.16
Nodes (23): CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath(), TestFindRoute_InexistentPathOnlyRoot(), TestFindRoute_InexistentRoot() (+15 more)

### Community 1 - "Tree"
Cohesion: 0.19
Nodes (13): HTTPMethods, Method, Node, RouterTree, Tree, isParam(), matchPath(), namedPath() (+5 more)

### Community 2 - "testing.T"
Cohesion: 0.24
Nodes (23): net/http/httptest.Server, testing.T, NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_Group(), TestRouter_MiddlewareOrder() (+15 more)

### Community 3 - "RouterContext"
Cohesion: 0.17
Nodes (10): net/http.Request, net/http.ResponseWriter, RouterContext, NewContext(), setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey(), TestContext_ValidKey() (+2 more)

### Community 4 - "Router"
Cohesion: 0.33
Nodes (5): net/http.HandlerFunc, sync.RWMutex, Router, _const.HTTPMethods, NewPrefixRouter()

### Community 5 - "net/http.Handler"
Cohesion: 0.36
Nodes (5): Chain, Middleware, Router, net/http.Handler, NewRouter()

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

- **Why does `CreateTree()` connect `CreateTree` to `Tree`, `testing.T`, `RouterContext`, `Router`?**
  _High betweenness centrality (0.198) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `CreateTree`?**
  _High betweenness centrality (0.163) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `CreateTree`, `Router`, `net/http.Handler`?**
  _High betweenness centrality (0.139) - this node is a cross-community bridge._
- **Are the 23 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestFindRoute()`) actually correct?**
  _`CreateTree()` has 23 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_Group()`) actually correct?**
  _`NewRouter()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/joaolaureano/go-router` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._