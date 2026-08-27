# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,351 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 152 nodes · 346 edges · 11 communities (8 shown, 3 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 50 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0acad4a0`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
- Node
- testing.T
- NewContext
- _const.HTTPMethods
- Router
- chain_test.go
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `Node` - 24 edges
3. `Tree` - 15 edges
4. `Router` - 14 edges
5. `setup()` - 11 edges
6. `HTTPMethods` - 11 edges
7. `NewContext()` - 11 edges
8. `newNode()` - 8 edges
9. `Chain` - 8 edges
10. `methodFilter` - 7 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `Node` --references--> `HTTPMethods`  [EXTRACTED]
  tree/node.go → const/const.go
- `CreateTree()` --calls--> `newNode()`  [INFERRED]
  tree/tree.go → tree/node.go

## Import Cycles
- None detected.

## Communities (11 total, 3 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.16
Nodes (23): CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath(), TestFindRoute_InexistentPathOnlyRoot(), TestFindRoute_InexistentRoot() (+15 more)

### Community 1 - "Node"
Cohesion: 0.18
Nodes (14): HTTPMethods, Method, methodFilter, Node, cloneMethod(), cloneNode(), filterByAnyMethod(), filterByMethod() (+6 more)

### Community 2 - "testing.T"
Cohesion: 0.17
Nodes (25): net/http/httptest.Server, testing.T, NewPrefixRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+17 more)

### Community 3 - "NewContext"
Cohesion: 0.13
Nodes (18): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+10 more)

### Community 5 - "Router"
Cohesion: 0.19
Nodes (10): Chain, Middleware, Router, github.com/joaolaureano/go-router/tree.RouterTree, net/http.Handler, net/http.HandlerFunc, sync.RWMutex, NewRouter() (+2 more)

### Community 6 - "chain_test.go"
Cohesion: 0.28
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.19
Nodes (10): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint(), isParam(), TestIsParam(), TestValidatePath_InvalidPaths() (+2 more)

## Knowledge Gaps
- **2 isolated node(s):** `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `CreateTree` to `Node`, `testing.T`, `NewContext`, `Router`, `newNode`?**
  _High betweenness centrality (0.163) - this node is a cross-community bridge._
- **Why does `Node` connect `Node` to `newNode`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `Tree` connect `Node` to `CreateTree`, `newNode`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewContext` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._