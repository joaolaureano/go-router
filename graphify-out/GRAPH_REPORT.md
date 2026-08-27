# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~6,338 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 170 nodes · 393 edges · 16 communities (7 shown, 9 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 76 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `55603c21`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
- node
- testing.T
- context.RouterContext
- _const.HTTPMethods
- methodFilter
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- Go-Router
- Node
- Tree
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 30 edges
2. `CreateTree()` - 30 edges
3. `node` - 18 edges
4. `Tree` - 16 edges
5. `Router` - 14 edges
6. `Method` - 14 edges
7. `setup()` - 11 edges
8. `assertFound()` - 9 edges
9. `Chain` - 9 edges
10. `newNode()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `TestNewRouterWithPrefix()` --calls--> `NewPrefixRouter()`  [INFERRED]
  router/router_test.go → router/router.go

## Import Cycles
- None detected.

## Communities (16 total, 9 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.14
Nodes (28): CreateTree(), assertFound(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_DiscardsValuesFromAbandonedBranches(), TestLookup_EmptyTree() (+20 more)

### Community 1 - "node"
Cohesion: 0.19
Nodes (11): endpoint, lookup, node, cloneEndpoint(), cloneNode(), mergeNodes(), newNode(), TestNewNodeInitializesInvariantState() (+3 more)

### Community 2 - "testing.T"
Cohesion: 0.17
Nodes (33): net/http/httptest.Server, testing.T, NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+25 more)

### Community 3 - "context.RouterContext"
Cohesion: 0.17
Nodes (14): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+6 more)

### Community 6 - "Router"
Cohesion: 0.11
Nodes (17): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+9 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.24
Nodes (12): Match, Method, Param, Status, Tree, isParam(), splitSegments(), TestIsParam() (+4 more)

## Knowledge Gaps
- **5 isolated node(s):** `Install`, `Instalação`, `routeTree`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `CreateTree` to `node`, `testing.T`, `Tree`, `Router`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `CreateTree`, `Router`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `Tree`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `routeTree` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CreateTree` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._