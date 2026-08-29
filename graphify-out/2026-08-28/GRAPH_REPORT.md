# Graph Report - go-router  (2026-08-28)

## Corpus Check
- 30 files · ~19,291 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 356 nodes · 1039 edges · 17 communities (13 shown, 4 thin omitted)
- Extraction: 74% EXTRACTED · 26% INFERRED · 0% AMBIGUOUS · INFERRED: 271 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3c9c3825`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- routing_test.go
- attack
- testing.T
- FromRequest
- Load tests
- run.sh script
- Router
- Repository
- github.com/joaolaureano/go-router
- Go-Router
- NewChain
- Go-Router
- bench_test.go
- node
- E

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 62 edges
2. `CreateTree()` - 48 edges
3. `Router` - 28 edges
4. `attack()` - 21 edges
5. `requireNoTransportErrors()` - 19 edges
6. `node` - 19 edges
7. `server()` - 18 edges
8. `requireLatencyBudget()` - 18 edges
9. `roundRobin()` - 18 edges
10. `skipShort()` - 18 edges

## Surprising Connections (you probably didn't know these)
- `TestRouter_VerbShortcuts()` --calls--> `echo()`  [INFERRED]
  router/router_test.go → loadtest/app/app.go
- `TestRouter_WithInheritsParentMiddleware()` --calls--> `tag()`  [INFERRED]
  router/router_test.go → loadtest/app/app.go
- `Router` --references--> `Middleware`  [EXTRACTED]
  router/router.go → chain/chain.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `TestNewRouterSatisfiesTheAcceptSurface()` --calls--> `NewRouter()`  [INFERRED]
  gorouter_test.go → gorouter.go

## Import Cycles
- None detected.

## Communities (17 total, 4 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.09
Nodes (48): CreateTree(), assertFound(), assertPanicsWith(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_CatchAll() (+40 more)

### Community 1 - "attack"
Cohesion: 0.20
Nodes (40): github.com/tsenart/vegeta/v12/lib.Metrics, github.com/tsenart/vegeta/v12/lib.Target, github.com/tsenart/vegeta/v12/lib.Targeter, net/http.Header, net/http/httptest.Server, time.Duration, TestLoad_MountWhileServing(), TestLoad_RegistrationWhileServing() (+32 more)

### Community 2 - "testing.T"
Cohesion: 0.08
Nodes (67): testing.T, Param(), TestParam(), NewPrefixRouter(), NewRouter(), setup(), TestDecodeSegmentLeavesMalformedEscapingAlone(), TestNewRouter() (+59 more)

### Community 3 - "FromRequest"
Cohesion: 0.14
Nodes (20): contextKey, routeContext, RouterContext, context.Context, net/http.Request, FromParams(), FromRequest(), NewContext() (+12 more)

### Community 4 - "Load tests"
Cohesion: 0.33
Nodes (5): Load tests, Running, The vegeta CLI, What each scenario pins, Why a separate module

### Community 6 - "Router"
Cohesion: 0.08
Nodes (26): Chain, Middleware, Router, net/http.Handler, net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer (+18 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "NewChain"
Cohesion: 0.24
Nodes (8): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestChainSealsOnceItIsBaked(), TestNewChain(), TestNewChain_Middlewares()

### Community 18 - "bench_test.go"
Cohesion: 0.40
Nodes (14): testing.B, BenchmarkCatchAll(), BenchmarkEscaped(), BenchmarkFanout50(), BenchmarkMethodNotAllowed(), BenchmarkMiss(), BenchmarkOptions(), BenchmarkParam() (+6 more)

### Community 20 - "node"
Cohesion: 0.23
Nodes (14): advertisedMethods(), endpoint, lookup, Method, node, cloneEndpoint(), cloneNode(), node[E] (+6 more)

### Community 28 - "E"
Cohesion: 0.11
Nodes (28): net/http/httptest.ResponseRecorder, testing.F, assertAllowHeader(), FuzzHeadFollowsGet(), fuzzRouter(), FuzzServeAnyTarget(), Router, serve() (+20 more)

## Knowledge Gaps
- **10 isolated node(s):** `github.com/joaolaureano/go-router`, `github.com/joaolaureano/go-router/loadtest`, `run.sh script`, `contextKey`, `Install` (+5 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewRouter()` connect `testing.T` to `attack`, `bench_test.go`, `E`, `Router`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `assertFound()` connect `routing_test.go` to `testing.T`, `node`, `E`, `Router`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `testing.T`, `NewChain`, `node`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Are the 56 inferred relationships involving `NewRouter()` (e.g. with `BenchmarkRegister300()` and `benchRouter()`) actually correct?**
  _`NewRouter()` has 56 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **Are the 16 inferred relationships involving `attack()` (e.g. with `TestLoad_MountWhileServing()` and `TestLoad_RegistrationWhileServing()`) actually correct?**
  _`attack()` has 16 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/joaolaureano/go-router`, `github.com/joaolaureano/go-router/loadtest`, `run.sh script` to the rest of the system?**
  _10 weakly-connected nodes found - possible documentation gaps or missing edges._