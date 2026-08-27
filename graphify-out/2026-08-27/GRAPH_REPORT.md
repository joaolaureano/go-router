# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~7,001 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 181 nodes · 424 edges · 17 communities (8 shown, 9 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 83 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `82d08629`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- tree_test.go
- node
- testing.T
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
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `CreateTree()` --calls--> `newNode()`  [INFERRED]
  tree/tree.go → tree/node.go
- `TestCreateTree()` --calls--> `CreateTree()`  [INFERRED]
  tree/tree_test.go → tree/tree.go

## Import Cycles
- None detected.

## Communities (17 total, 9 thin omitted)

### Community 0 - "tree_test.go"
Cohesion: 0.14
Nodes (28): CreateTree(), assertFound(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_DiscardsValuesFromAbandonedBranches(), TestLookup_EmptyTree() (+20 more)

### Community 1 - "node"
Cohesion: 0.28
Nodes (7): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "testing.T"
Cohesion: 0.14
Nodes (39): net/http/httptest.Server, testing.T, NewPrefixRouter(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod() (+31 more)

### Community 3 - "FromRequest"
Cohesion: 0.17
Nodes (17): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param() (+9 more)

### Community 6 - "router.Router"
Cohesion: 0.09
Nodes (19): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+11 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.26
Nodes (11): Match, Param, Status, Tree, isParam(), splitSegments(), TestIsParam(), TestValidatePath_InvalidPaths() (+3 more)

### Community 16 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

## Knowledge Gaps
- **5 isolated node(s):** `Install`, `Instalação`, `contextKey`, `routeTree`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `tree_test.go` to `newNode`, `testing.T`, `Tree`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `testing.T` to `router.Router`?**
  _High betweenness centrality (0.080) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **Are the 32 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 32 INFERRED edges - model-reasoned connections that need verification._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tree_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._