# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~6,370 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 170 nodes · 371 edges · 17 communities (8 shown, 9 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 51 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `79b7a745`
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
- newNode

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `node` - 18 edges
3. `Tree` - 16 edges
4. `Method` - 14 edges
5. `Router` - 14 edges
6. `setup()` - 11 edges
7. `FromRequest()` - 10 edges
8. `assertFound()` - 9 edges
9. `Chain` - 9 edges
10. `newNode()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `CreateTree()` --calls--> `newNode()`  [INFERRED]
  tree/tree.go → tree/node.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go

## Import Cycles
- None detected.

## Communities (17 total, 9 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.14
Nodes (28): CreateTree(), assertFound(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_DiscardsValuesFromAbandonedBranches(), TestLookup_EmptyTree() (+20 more)

### Community 1 - "node"
Cohesion: 0.26
Nodes (7): endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), mergeNodes()

### Community 2 - "testing.T"
Cohesion: 0.12
Nodes (33): net/http/httptest.Server, testing.T, NewPrefixRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+25 more)

### Community 3 - "FromRequest"
Cohesion: 0.19
Nodes (15): contextKey, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), Param(), setup() (+7 more)

### Community 6 - "Router"
Cohesion: 0.11
Nodes (17): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+9 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.29
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

- **Why does `CreateTree()` connect `CreateTree` to `newNode`, `testing.T`, `Tree`, `Router`?**
  _High betweenness centrality (0.091) - this node is a cross-community bridge._
- **Why does `Method` connect `node` to `CreateTree`, `FromRequest`, `Tree`, `Router`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CreateTree` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.12477718360071301 - nodes in this community are weakly interconnected._