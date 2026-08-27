# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 18 files · ~8,846 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 217 nodes · 516 edges · 21 communities (11 shown, 10 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 99 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d86bac58`
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
- net/http.Handler
- Tree
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- Method
- NewChain

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 43 edges
2. `CreateTree()` - 36 edges
3. `Router` - 24 edges
4. `Tree` - 20 edges
5. `node` - 20 edges
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
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go

## Import Cycles
- None detected.

## Communities (21 total, 10 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.14
Nodes (39): testing.T, CreateTree(), assertFound(), assertPanicsWith(), Method, TestCreateTree(), TestIsParam(), TestLookup() (+31 more)

### Community 1 - "node"
Cohesion: 0.17
Nodes (14): endpoint, lookup, node, cloneEndpoint(), cloneNode(), Method, mergeBranch(), mergeNodes() (+6 more)

### Community 2 - "NewRouter"
Cohesion: 0.10
Nodes (41): net/http/httptest.Server, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_CatchAllRoute(), TestRouter_CatchAllUnderGroup(), TestRouter_Group() (+33 more)

### Community 3 - "FromRequest"
Cohesion: 0.16
Nodes (17): contextKey, param, RouterContext, net/http.Request, FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+9 more)

### Community 6 - "Router"
Cohesion: 0.22
Nodes (7): net/http.HandlerFunc, net/http.ResponseWriter, sync.RWMutex, Router, defaultMethodNotAllowed(), Method, Method

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "net/http.Handler"
Cohesion: 0.17
Nodes (9): Chain, Middleware, Router, net/http.Handler, NewPrefixRouter(), NewRouter(), NewPrefixRouter(), TestNewRouterWithPrefix() (+1 more)

### Community 12 - "Tree"
Cohesion: 0.22
Nodes (14): node, segmentKind, Match, Param, Status, classify(), Method, Tree (+6 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Router` connect `Router` to `NewRouter`, `NewChain`, `net/http.Handler`, `Tree`?**
  _High betweenness centrality (0.176) - this node is a cross-community bridge._
- **Why does `Tree` connect `Tree` to `testing.T`, `node`, `Router`?**
  _High betweenness centrality (0.129) - this node is a cross-community bridge._
- **Why does `CreateTree()` connect `testing.T` to `node`, `net/http.Handler`, `Tree`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Are the 40 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_BacktracksWhenStaticMatchLacksMethod()`) actually correct?**
  _`NewRouter()` has 40 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `CreateTree()` (e.g. with `newNode()` and `TestCreateTree()`) actually correct?**
  _`CreateTree()` has 34 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.14487179487179488 - nodes in this community are weakly interconnected._