# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,787 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 164 nodes · 389 edges · 11 communities (8 shown, 3 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 77 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f6aa4881`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
- Node
- testing.T
- NewContext
- _const.HTTPMethods
- Tree
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `NewRouter()` - 28 edges
3. `Node` - 20 edges
4. `Tree` - 16 edges
5. `Router` - 14 edges
6. `NewContext()` - 11 edges
7. `HTTPMethods` - 11 edges
8. `setup()` - 11 edges
9. `Chain` - 9 edges
10. `RouterContext` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `TestFindRoute_FallsBackToParameterAfterStaticBranchFails()` --calls--> `NewContext()`  [EXTRACTED]
  tree/tree_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (11 total, 3 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.16
Nodes (23): CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath(), TestFindRoute_InexistentPathOnlyRoot(), TestFindRoute_InexistentRoot() (+15 more)

### Community 1 - "Node"
Cohesion: 0.22
Nodes (11): HTTPMethods, Method, methodFilter, Node, cloneMethod(), cloneNode(), filterByAnyMethod(), filterByMethod() (+3 more)

### Community 2 - "testing.T"
Cohesion: 0.18
Nodes (31): net/http/httptest.Server, testing.T, NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+23 more)

### Community 3 - "NewContext"
Cohesion: 0.13
Nodes (18): contextKey, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+10 more)

### Community 5 - "Tree"
Cohesion: 0.21
Nodes (12): github.com/joaolaureano/go-router/router/context.RouterContext, methodFilter, Node, tree.RouterTree, Tree, isParam(), setPathVariableValues(), splitSegments() (+4 more)

### Community 6 - "Router"
Cohesion: 0.11
Nodes (17): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+9 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

## Knowledge Gaps
- **2 isolated node(s):** `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `CreateTree` to `testing.T`, `NewContext`, `Tree`, `Router`, `newNode`?**
  _High betweenness centrality (0.165) - this node is a cross-community bridge._
- **Why does `Node` connect `Node` to `newNode`?**
  _High betweenness centrality (0.113) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `CreateTree`, `Router`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewContext` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._