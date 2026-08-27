# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~6,170 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 165 nodes · 387 edges · 15 communities (8 shown, 7 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 76 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f977b01d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- node
- NewRouter
- context.RouterContext
- _const.HTTPMethods
- methodFilter
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode
- Node
- Tree
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `NewRouter()` - 30 edges
3. `node` - 18 edges
4. `Tree` - 16 edges
5. `Router` - 14 edges
6. `HTTPMethods` - 11 edges
7. `setup()` - 11 edges
8. `assertFound()` - 9 edges
9. `Chain` - 9 edges
10. `newNode()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `assertFound()` --references--> `HTTPMethods`  [EXTRACTED]
  tree/tree_test.go → const/const.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go

## Import Cycles
- None detected.

## Communities (15 total, 7 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.19
Nodes (30): testing.T, CreateTree(), assertFound(), TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+22 more)

### Community 1 - "node"
Cohesion: 0.24
Nodes (7): HTTPMethods, endpoint, lookup, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "NewRouter"
Cohesion: 0.14
Nodes (30): net/http/httptest.Server, NewRouter(), setup(), TestNewRouter(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group(), TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute(), TestRouter_GroupNotFoundReachesServingRouter() (+22 more)

### Community 3 - "context.RouterContext"
Cohesion: 0.17
Nodes (14): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+6 more)

### Community 6 - "Router"
Cohesion: 0.10
Nodes (19): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+11 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

### Community 12 - "Tree"
Cohesion: 0.29
Nodes (10): Match, Param, Status, Tree, isParam(), splitSegments(), TestValidatePath_InvalidPaths(), TestValidatePath_ValidPaths() (+2 more)

## Knowledge Gaps
- **3 isolated node(s):** `contextKey`, `routeTree`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `testing.T` to `NewRouter`, `newNode`, `Tree`, `Router`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `testing.T`, `Router`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.091) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `routeTree`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewRouter` be split into smaller, more focused modules?**
  _Cohesion score 0.13763440860215054 - nodes in this community are weakly interconnected._