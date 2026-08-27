# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,613 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 159 nodes · 378 edges · 11 communities (8 shown, 3 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 73 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `23d52bf3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- Node
- NewRouter
- NewContext
- _const.HTTPMethods
- net/http.Handler
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `NewRouter()` - 26 edges
3. `Node` - 25 edges
4. `Tree` - 15 edges
5. `Router` - 14 edges
6. `setup()` - 11 edges
7. `HTTPMethods` - 11 edges
8. `NewContext()` - 11 edges
9. `Chain` - 9 edges
10. `newNode()` - 8 edges

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

## Communities (11 total, 3 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.17
Nodes (31): testing.T, CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_FallsBackToParameterAfterStaticBranchFails(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath() (+23 more)

### Community 1 - "Node"
Cohesion: 0.15
Nodes (16): HTTPMethods, github.com/joaolaureano/go-router/router/context.RouterContext, methodFilter, Node, cloneMethod(), cloneNode(), filterByAnyMethod(), filterByMethod() (+8 more)

### Community 2 - "NewRouter"
Cohesion: 0.16
Nodes (26): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group(), TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute(), TestRouter_GroupRejectsMiddlewareAfterItsFirstRoute() (+18 more)

### Community 3 - "NewContext"
Cohesion: 0.17
Nodes (14): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+6 more)

### Community 5 - "net/http.Handler"
Cohesion: 0.24
Nodes (6): Chain, Middleware, Router, net/http.Handler, NewRouter(), Method

### Community 6 - "Router"
Cohesion: 0.14
Nodes (14): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares(), github.com/joaolaureano/go-router/tree.RouterTree (+6 more)

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

- **Why does `CreateTree()` connect `testing.T` to `Node`, `NewRouter`, `newNode`, `Router`?**
  _High betweenness centrality (0.166) - this node is a cross-community bridge._
- **Why does `Node` connect `Node` to `newNode`, `net/http.Handler`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `net/http.Handler`, `Router`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 23 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 23 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Router` be split into smaller, more focused modules?**
  _Cohesion score 0.14035087719298245 - nodes in this community are weakly interconnected._