# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 16 files · ~8,526 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 212 nodes · 511 edges · 20 communities (9 shown, 11 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 98 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `50e4155a`
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
- github.com/joaolaureano/go-router/router.Router
- Method
- NewChain
- newNode

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 43 edges
2. `CreateTree()` - 35 edges
3. `Router` - 26 edges
4. `node` - 22 edges
5. `Tree` - 20 edges
6. `setup()` - 11 edges
7. `FromRequest()` - 11 edges
8. `Param()` - 10 edges
9. `assertFound()` - 9 edges
10. `Chain` - 9 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `Router` --references--> `Tree`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `NewRouter()` --references--> `Router`  [EXTRACTED]
  gorouter.go → router/router.go

## Import Cycles
- None detected.

## Communities (20 total, 11 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.15
Nodes (37): testing.T, CreateTree(), assertFound(), Method, TestCreateTree(), TestIsParam(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath() (+29 more)

### Community 1 - "node"
Cohesion: 0.24
Nodes (9): endpoint, lookup, node, cloneEndpoint(), cloneNode(), Method, mergeBranch(), mergeNodes() (+1 more)

### Community 2 - "NewRouter"
Cohesion: 0.10
Nodes (41): net/http/httptest.Server, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_CatchAllRoute(), TestRouter_CatchAllUnderGroup(), TestRouter_Group() (+33 more)

### Community 3 - "FromRequest"
Cohesion: 0.13
Nodes (20): contextKey, param, RouterContext, net/http.Request, net/http.ResponseWriter, FromRequest(), NewContext(), setup() (+12 more)

### Community 6 - "Router"
Cohesion: 0.18
Nodes (10): Router, net/http.HandlerFunc, sync.RWMutex, NewPrefixRouter(), NewRouter(), Method, Router, NewPrefixRouter() (+2 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 12 - "Tree"
Cohesion: 0.16
Nodes (15): Chain, Middleware, net/http.Handler, Match, Param, Status, classify(), Method (+7 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 19 - "newNode"
Cohesion: 0.53
Nodes (5): newNode(), TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint()

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Router` connect `Router` to `NewRouter`, `NewChain`, `FromRequest`, `Tree`?**
  _High betweenness centrality (0.195) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `testing.T`, `node`, `Router`?**
  _High betweenness centrality (0.168) - this node is a cross-community bridge._
- **Why does `node` connect `node` to `newNode`, `Tree`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Are the 40 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 40 INFERRED edges - model-reasoned connections that need verification._
- **Are the 33 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `NewRouter` be split into smaller, more focused modules?**
  _Cohesion score 0.10336817653890824 - nodes in this community are weakly interconnected._