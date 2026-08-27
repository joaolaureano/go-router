# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 18 files · ~10,003 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 239 nodes · 581 edges · 24 communities (11 shown, 13 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 102 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cf9a8950`
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
- Chain
- Tree
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- NewChain
- E
- Method
- Method
- Method

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 52 edges
2. `CreateTree()` - 36 edges
3. `Router` - 26 edges
4. `node` - 17 edges
5. `Tree` - 17 edges
6. `Param()` - 12 edges
7. `FromRequest()` - 11 edges
8. `setup()` - 11 edges
9. `assertFound()` - 10 edges
10. `RouterContext` - 9 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → tree/tree.go
- `requestSegments()` --calls--> `SplitPath()`  [EXTRACTED]
  router/router.go → tree/tree.go

## Import Cycles
- None detected.

## Communities (24 total, 13 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.12
Nodes (43): testing.T, TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint(), CreateTree(), assertFound(), assertPanicsWith() (+35 more)

### Community 1 - "node"
Cohesion: 0.28
Nodes (12): endpoint, lookup, tree.Method, node, cloneEndpoint(), cloneNode(), node[E], E (+4 more)

### Community 2 - "NewRouter"
Cohesion: 0.09
Nodes (50): net/http/httptest.Server, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_AnswersOptionsAutomatically(), TestRouter_CatchAllRoute(), TestRouter_CatchAllUnderGroup() (+42 more)

### Community 3 - "FromRequest"
Cohesion: 0.17
Nodes (17): contextKey, RouterContext, github.com/joaolaureano/go-router/tree.Param, FromParams(), FromRequest(), NewContext(), setup(), TestContext_InvalidKey() (+9 more)

### Community 6 - "Router"
Cohesion: 0.17
Nodes (11): github.com/joaolaureano/go-router/tree.Method, net/http.HandlerFunc, net/http.Request, net/http.ResponseWriter, sync.RWMutex, Method, Router, advertisedMethods() (+3 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.17
Nodes (9): Chain, Middleware, Router, net/http.Handler, NewPrefixRouter(), NewRouter(), NewPrefixRouter(), TestNewRouterWithPrefix() (+1 more)

### Community 12 - "Tree"
Cohesion: 0.21
Nodes (17): E, node, segmentKind, Match, tree.Param, Status, Tree, classify() (+9 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

## Knowledge Gaps
- **4 isolated node(s):** `contextKey`, `Install`, `Instalação`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `assertFound()` connect `testing.T` to `node`, `Chain`, `Tree`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewRouter`, `NewChain`, `Chain`?**
  _High betweenness centrality (0.157) - this node is a cross-community bridge._
- **Why does `NewRouter()` connect `NewRouter` to `FromRequest`, `Chain`, `Router`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Are the 49 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_AnswersOptionsAutomatically()`) actually correct?**
  _`NewRouter()` has 49 INFERRED edges - model-reasoned connections that need verification._
- **Are the 33 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `Install`, `Instalação` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.12323232323232323 - nodes in this community are weakly interconnected._