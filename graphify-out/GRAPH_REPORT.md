# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~7,001 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 180 nodes · 424 edges · 18 communities (9 shown, 9 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 83 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5d6f33e0`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- node
- NewRouter
- FromRequest
- _const.HTTPMethods
- methodFilter
- router.Router
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
- NewChain

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 35 edges
2. `CreateTree()` - 29 edges
3. `node` - 18 edges
4. `Tree` - 16 edges
5. `Method` - 13 edges
6. `setup()` - 11 edges
7. `FromRequest()` - 11 edges
8. `assertFound()` - 9 edges
9. `Chain` - 9 edges
10. `newNode()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `CreateTree()` --calls--> `newNode()`  [INFERRED]
  tree/tree.go → tree/node.go
- `TestIsParam()` --calls--> `isParam()`  [INFERRED]
  tree/tree_test.go → tree/tree.go

## Import Cycles
- None detected.

## Communities (18 total, 9 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.19
Nodes (30): testing.T, CreateTree(), assertFound(), TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+22 more)

### Community 1 - "node"
Cohesion: 0.26
Nodes (7): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "NewRouter"
Cohesion: 0.13
Nodes (30): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_Group(), TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute(), TestRouter_GroupMethodNotAllowedReachesServingRouter(), TestRouter_GroupNotFoundReachesServingRouter() (+22 more)

### Community 3 - "FromRequest"
Cohesion: 0.13
Nodes (22): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param() (+14 more)

### Community 6 - "router.Router"
Cohesion: 0.13
Nodes (14): Chain, Middleware, Router, net/http.Handler, net/http.HandlerFunc, sync.RWMutex, NewPrefixRouter(), NewRouter() (+6 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.32
Nodes (10): Match, Param, Status, Tree, isParam(), splitSegments(), TestValidatePath_InvalidPaths(), TestValidatePath_ValidPaths() (+2 more)

### Community 16 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

### Community 17 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

## Knowledge Gaps
- **5 isolated node(s):** `Install`, `Instalação`, `contextKey`, `routeTree`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `newNode`, `Tree`, `router.Router`?**
  _High betweenness centrality (0.078) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **Why does `Method` connect `node` to `testing.T`, `FromRequest`, `Tree`?**
  _High betweenness centrality (0.073) - this node is a cross-community bridge._
- **Are the 32 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 32 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewRouter` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._