# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~6,678 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 174 nodes · 407 edges · 18 communities (9 shown, 9 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 80 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7f62e744`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- node
- NewRouter
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
- newNode
- net/http.Handler

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
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `CreateTree()` --calls--> `newNode()`  [INFERRED]
  tree/tree.go → tree/node.go

## Import Cycles
- None detected.

## Communities (18 total, 9 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.19
Nodes (30): testing.T, CreateTree(), assertFound(), TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+22 more)

### Community 1 - "node"
Cohesion: 0.28
Nodes (7): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "NewRouter"
Cohesion: 0.13
Nodes (32): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group(), TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute(), TestRouter_GroupNotFoundReachesServingRouter() (+24 more)

### Community 3 - "FromRequest"
Cohesion: 0.18
Nodes (16): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param() (+8 more)

### Community 6 - "Router"
Cohesion: 0.12
Nodes (16): Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares() (+8 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.29
Nodes (10): Match, Param, Status, Tree, isParam(), splitSegments(), TestValidatePath_InvalidPaths(), TestValidatePath_ValidPaths() (+2 more)

### Community 16 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

### Community 17 - "net/http.Handler"
Cohesion: 0.36
Nodes (4): Chain, Router, net/http.Handler, NewRouter()

## Knowledge Gaps
- **5 isolated node(s):** `contextKey`, `Install`, `Instalação`, `routeTree`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `net/http.Handler`, `Router`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewRouter`, `FromRequest`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Why does `CreateTree()` connect `testing.T` to `newNode`, `NewRouter`, `Tree`, `Router`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Are the 29 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `Install`, `Instalação` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewRouter` be split into smaller, more focused modules?**
  _Cohesion score 0.12878787878787878 - nodes in this community are weakly interconnected._