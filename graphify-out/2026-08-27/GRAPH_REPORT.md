# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~6,597 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 173 nodes · 405 edges · 17 communities (7 shown, 10 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 80 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cc300fd2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
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
- Router

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 32 edges
2. `CreateTree()` - 30 edges
3. `node` - 18 edges
4. `Tree` - 16 edges
5. `Router` - 14 edges
6. `Method` - 13 edges
7. `FromRequest()` - 11 edges
8. `setup()` - 11 edges
9. `assertFound()` - 9 edges
10. `Chain` - 9 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `TestNewRouterWithPrefix()` --calls--> `NewPrefixRouter()`  [INFERRED]
  router/router_test.go → router/router.go

## Import Cycles
- None detected.

## Communities (17 total, 10 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.14
Nodes (28): CreateTree(), assertFound(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_DiscardsValuesFromAbandonedBranches(), TestLookup_EmptyTree() (+20 more)

### Community 1 - "node"
Cohesion: 0.18
Nodes (11): endpoint, lookup, node, cloneEndpoint(), cloneNode(), mergeNodes(), newNode(), TestNewNodeInitializesInvariantState() (+3 more)

### Community 2 - "testing.T"
Cohesion: 0.16
Nodes (35): net/http/httptest.Server, testing.T, NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+27 more)

### Community 3 - "FromRequest"
Cohesion: 0.19
Nodes (15): contextKey, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param(), setup() (+7 more)

### Community 6 - "Router"
Cohesion: 0.13
Nodes (14): Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares() (+6 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.19
Nodes (14): Chain, net/http.Handler, Match, Method, Param, Status, Tree, isParam() (+6 more)

## Knowledge Gaps
- **5 isolated node(s):** `contextKey`, `routeTree`, `Install`, `Instalação`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewRouter()` connect `testing.T` to `CreateTree`, `Router`, `Router`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Why does `CreateTree()` connect `CreateTree` to `node`, `testing.T`, `Tree`, `Router`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `testing.T`, `FromRequest`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Are the 29 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `routeTree`, `Install` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CreateTree` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._