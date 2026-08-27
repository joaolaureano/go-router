# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~6,931 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 179 nodes · 421 edges · 17 communities (8 shown, 9 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 83 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `128e021a`
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
- newNode

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 35 edges
2. `CreateTree()` - 29 edges
3. `node` - 18 edges
4. `Router` - 16 edges
5. `Tree` - 16 edges
6. `Method` - 13 edges
7. `setup()` - 11 edges
8. `FromRequest()` - 11 edges
9. `assertFound()` - 9 edges
10. `Chain` - 9 edges

## Surprising Connections (you probably didn't know these)
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `CreateTree()` --calls--> `newNode()`  [INFERRED]
  tree/tree.go → tree/node.go
- `TestIsParam()` --calls--> `isParam()`  [INFERRED]
  tree/tree_test.go → tree/tree.go

## Import Cycles
- None detected.

## Communities (17 total, 9 thin omitted)

### Community 0 - "tree_test.go"
Cohesion: 0.14
Nodes (29): CreateTree(), assertFound(), TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_DiscardsValuesFromAbandonedBranches() (+21 more)

### Community 1 - "node"
Cohesion: 0.26
Nodes (7): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "testing.T"
Cohesion: 0.14
Nodes (39): net/http/httptest.Server, testing.T, NewPrefixRouter(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod() (+31 more)

### Community 3 - "FromRequest"
Cohesion: 0.17
Nodes (17): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param() (+9 more)

### Community 6 - "Router"
Cohesion: 0.10
Nodes (17): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+9 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.32
Nodes (10): Match, Param, Status, Tree, isParam(), splitSegments(), TestValidatePath_InvalidPaths(), TestValidatePath_ValidPaths() (+2 more)

### Community 16 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

## Knowledge Gaps
- **5 isolated node(s):** `Install`, `Instalação`, `routeTree`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Router` connect `Router` to `testing.T`, `FromRequest`?**
  _High betweenness centrality (0.110) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `Router`?**
  _High betweenness centrality (0.080) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Are the 32 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 32 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `routeTree` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tree_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.135632183908046 - nodes in this community are weakly interconnected._