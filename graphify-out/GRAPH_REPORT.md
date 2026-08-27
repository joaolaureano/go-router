# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~5,877 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 165 nodes · 368 edges · 12 communities (7 shown, 5 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 52 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `727520fb`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CreateTree
- Node
- testing.T
- NewContext
- _const.HTTPMethods
- methodFilter
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- newNode
- Node

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 30 edges
2. `Node` - 25 edges
3. `Tree` - 16 edges
4. `Router` - 14 edges
5. `setup()` - 11 edges
6. `HTTPMethods` - 11 edges
7. `NewContext()` - 11 edges
8. `Chain` - 9 edges
9. `newNode()` - 8 edges
10. `FromRequest()` - 8 edges

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

## Communities (12 total, 5 thin omitted)

### Community 0 - "CreateTree"
Cohesion: 0.16
Nodes (23): CreateTree(), TestCreateTree(), TestFindRoute(), TestFindRoute_EmptyPath(), TestFindRoute_InexistentMethod(), TestFindRoute_InexistentPath(), TestFindRoute_InexistentPathOnlyRoot(), TestFindRoute_InexistentRoot() (+15 more)

### Community 1 - "Node"
Cohesion: 0.13
Nodes (20): HTTPMethods, github.com/joaolaureano/go-router/router/context.RouterContext, Method, methodFilter, Node, cloneMethod(), cloneNode(), filterByAnyMethod() (+12 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (32): net/http/httptest.Server, testing.T, NewPrefixRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_BacktracksWhenStaticMatchLacksMethod(), TestRouter_Group() (+24 more)

### Community 3 - "NewContext"
Cohesion: 0.13
Nodes (18): contextKey, context.RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+10 more)

### Community 6 - "Router"
Cohesion: 0.11
Nodes (17): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+9 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 10 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesStaticAndParameterChildren(), TestNodeSetEndpoint()

## Knowledge Gaps
- **2 isolated node(s):** `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `CreateTree` to `Node`, `testing.T`, `NewContext`, `Router`, `newNode`?**
  _High betweenness centrality (0.149) - this node is a cross-community bridge._
- **Why does `Node` connect `Node` to `newNode`?**
  _High betweenness centrality (0.130) - this node is a cross-community bridge._
- **Why does `Tree` connect `Node` to `CreateTree`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `github.com/joaolaureano/go-router` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Node` be split into smaller, more focused modules?**
  _Cohesion score 0.12857142857142856 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.12878787878787878 - nodes in this community are weakly interconnected._