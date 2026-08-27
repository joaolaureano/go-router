# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 19 files · ~11,264 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 272 nodes · 634 edges · 35 communities (11 shown, 24 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 112 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9de3428f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- routing_test.go
- tree.Method
- testing.T
- FromRequest
- _const.HTTPMethods
- methodFilter
- Router
- Repository
- github.com/joaolaureano/go-router
- _const.HTTPMethods
- Go-Router
- Chain
- sync.RWMutex
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- E
- E
- node
- E
- node
- node_test.go
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- Method
- segmentKind
- Param
- E
- tree.Param
- Method

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 55 edges
2. `CreateTree()` - 38 edges
3. `Router` - 24 edges
4. `node` - 20 edges
5. `Param()` - 14 edges
6. `FromRequest()` - 12 edges
7. `state` - 12 edges
8. `assertFound()` - 11 edges
9. `setup()` - 11 edges
10. `RouterContext` - 10 edges

## Surprising Connections (you probably didn't know these)
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `TestParam()` --calls--> `Param()`  [INFERRED]
  router/context/context_test.go → router/context/context.go

## Import Cycles
- None detected.

## Communities (35 total, 24 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.10
Nodes (39): CreateTree(), assertFound(), assertPanicsWith(), Method, TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters() (+31 more)

### Community 2 - "testing.T"
Cohesion: 0.10
Nodes (59): net/http/httptest.Server, testing.T, Param(), NewRouter(), setup(), TestNewRouter(), TestNewRouterWithPrefix(), TestRouter_AnswersOptionsAutomatically() (+51 more)

### Community 3 - "FromRequest"
Cohesion: 0.14
Nodes (20): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, FromParams(), FromRequest() (+12 more)

### Community 6 - "Router"
Cohesion: 0.12
Nodes (15): net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer, sync.Mutex, Method, Router, advertisedMethods() (+7 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.12
Nodes (13): Chain, Middleware, NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain() (+5 more)

### Community 17 - "E"
Cohesion: 0.20
Nodes (15): Match, routing.Param, classify(), E, Method, isParam(), nameParams(), splitSegments() (+7 more)

### Community 20 - "node"
Cohesion: 0.24
Nodes (13): endpoint, lookup, node, cloneEndpoint(), cloneNode(), node[E], E, Method (+5 more)

### Community 23 - "node_test.go"
Cohesion: 0.40
Nodes (4): TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint()

## Knowledge Gaps
- **4 isolated node(s):** `contextKey`, `Install`, `Instalação`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `routing_test.go` to `E`, `Router`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `Tree` connect `E` to `routing_test.go`, `node`?**
  _High betweenness centrality (0.141) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `testing.T`, `Chain`?**
  _High betweenness centrality (0.135) - this node is a cross-community bridge._
- **Are the 52 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_AnswersOptionsAutomatically()`) actually correct?**
  _`NewRouter()` has 52 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 34 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `Install`, `Instalação` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `routing_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10256410256410256 - nodes in this community are weakly interconnected._