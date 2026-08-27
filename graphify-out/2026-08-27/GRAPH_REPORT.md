# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 14 files · ~4,919 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 130 nodes · 306 edges · 9 communities (7 shown, 2 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 52 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `95dd0a1a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
- Tree
- testing.T
- RouterContext
- Router
- chain_test.go
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 29 edges
2. `Tree` - 21 edges
3. `NewRouter()` - 21 edges
4. `Node` - 13 edges
5. `Router` - 13 edges
6. `setup()` - 11 edges
7. `RouterContext` - 8 edges
8. `NewContext()` - 8 edges
9. `Chain` - 8 edges
10. `mergeNodes()` - 6 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `setPathVariableValues()` --references--> `RouterContext`  [EXTRACTED]
  tree/tree.go → router/context/context.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go

## Import Cycles
- None detected.

## Communities (9 total, 2 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.14
Nodes (28): NewContext(), CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_FallsBackToParameterAfterStaticBranchFails(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath() (+20 more)

### Community 1 - "Tree"
Cohesion: 0.18
Nodes (16): HTTPMethods, Method, Node, RouterTree, Tree, cloneMethod(), cloneNode(), isParam() (+8 more)

### Community 2 - "testing.T"
Cohesion: 0.22
Nodes (24): net/http/httptest.Server, testing.T, NewPrefixRouter(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_Group() (+16 more)

### Community 3 - "RouterContext"
Cohesion: 0.18
Nodes (8): net/http.Request, net/http.ResponseWriter, RouterContext, setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey(), TestContext_ValidKey(), TestInjectIntoRequest()

### Community 5 - "Router"
Cohesion: 0.19
Nodes (9): Chain, Middleware, Router, net/http.Handler, net/http.HandlerFunc, sync.RWMutex, NewRouter(), Router (+1 more)

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

- **Why does `CreateTree()` connect `CreateTree` to `Tree`, `testing.T`?**
  _High betweenness centrality (0.219) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `CreateTree`?**
  _High betweenness centrality (0.189) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `CreateTree`, `Router`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Are the 26 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestFindRoute()`) actually correct?**
  _`CreateTree()` has 26 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_Group()`) actually correct?**
  _`NewRouter()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/joaolaureano/go-router` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CreateTree` be split into smaller, more focused modules?**
  _Cohesion score 0.13793103448275862 - nodes in this community are weakly interconnected._