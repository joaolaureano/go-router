# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~7,735 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 199 nodes · 470 edges · 19 communities (8 shown, 11 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 90 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `79ceb873`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- tree_test.go
- node
- testing.T
- FromRequest
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
- github.com/joaolaureano/go-router/router.Router
- Method
- NewChain

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 41 edges
2. `CreateTree()` - 29 edges
3. `Router` - 26 edges
4. `node` - 19 edges
5. `Tree` - 18 edges
6. `setup()` - 11 edges
7. `FromRequest()` - 11 edges
8. `Chain` - 9 edges
9. `newNode()` - 8 edges
10. `assertFound()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `Router` --references--> `Tree`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go

## Import Cycles
- None detected.

## Communities (19 total, 11 thin omitted)

### Community 0 - "tree_test.go"
Cohesion: 0.14
Nodes (28): CreateTree(), assertFound(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_DiscardsValuesFromAbandonedBranches(), TestLookup_EmptyTree() (+20 more)

### Community 1 - "node"
Cohesion: 0.19
Nodes (12): endpoint, lookup, node, cloneEndpoint(), cloneNode(), Method, mergeNodes(), newNode() (+4 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (45): net/http/httptest.Server, testing.T, NewPrefixRouter(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod() (+37 more)

### Community 3 - "FromRequest"
Cohesion: 0.16
Nodes (18): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param() (+10 more)

### Community 6 - "Router"
Cohesion: 0.15
Nodes (10): Chain, Middleware, Router, net/http.Handler, net/http.HandlerFunc, sync.RWMutex, NewPrefixRouter(), NewRouter() (+2 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.23
Nodes (12): Match, Param, Status, Method, Tree, isParam(), splitSegments(), TestIsParam() (+4 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Router` connect `Router` to `testing.T`, `NewChain`, `FromRequest`, `Tree`?**
  _High betweenness centrality (0.208) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `tree_test.go`, `node`, `Router`?**
  _High betweenness centrality (0.140) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `Router`?**
  _High betweenness centrality (0.109) - this node is a cross-community bridge._
- **Are the 38 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 38 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tree_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._