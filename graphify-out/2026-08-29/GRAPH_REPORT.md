# Graph Report - go-router  (2026-08-29)

## Corpus Check
- 33 files · ~22,501 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 380 nodes · 1119 edges · 18 communities (14 shown, 4 thin omitted)
- Extraction: 74% EXTRACTED · 26% INFERRED · 0% AMBIGUOUS · INFERRED: 290 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b379ed97`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- requireNoTransportErrors
- NewRouter
- context_test.go
- Load tests
- run.sh script
- Router
- Repository
- github.com/joaolaureano/go-router
- NewCatalog
- Go-Router
- NewChain
- Go-Router
- bench_test.go
- node
- routing/fuzz_test.go

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 63 edges
2. `CreateTree()` - 48 edges
3. `Router` - 29 edges
4. `requireNoTransportErrors()` - 23 edges
5. `skipShort()` - 22 edges
6. `attack()` - 21 edges
7. `requireStatus()` - 21 edges
8. `requireBody()` - 20 edges
9. `node` - 19 edges
10. `server()` - 18 edges

## Surprising Connections (you probably didn't know these)
- `TestRouter_VerbShortcuts()` --calls--> `echo()`  [INFERRED]
  router/router_test.go → loadtest/app/app.go
- `TestRouter_WithInheritsParentMiddleware()` --calls--> `tag()`  [INFERRED]
  router/router_test.go → loadtest/app/app.go
- `echo()` --calls--> `Param()`  [EXTRACTED]
  loadtest/app/app.go → router/context/context.go
- `New()` --calls--> `NewRouter()`  [EXTRACTED]
  loadtest/app/app.go → router/router.go
- `adminRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  loadtest/app/app.go → router/router.go

## Import Cycles
- None detected.

## Communities (18 total, 4 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.10
Nodes (55): testing.T, TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint(), TestNodeSetEndpointOnlyAppends(), CreateTree(), assertFound() (+47 more)

### Community 1 - "requireNoTransportErrors"
Cohesion: 0.16
Nodes (49): github.com/tsenart/vegeta/v12/lib.Attacker, github.com/tsenart/vegeta/v12/lib.Metrics, github.com/tsenart/vegeta/v12/lib.Pacer, github.com/tsenart/vegeta/v12/lib.Target, github.com/tsenart/vegeta/v12/lib.Targeter, net/http.Header, time.Duration, burst() (+41 more)

### Community 2 - "NewRouter"
Cohesion: 0.06
Nodes (64): contextKey, routeContext, context.Context, net/http/httptest.Server, net/http.Request, FromRequest(), Param(), TestWithParams() (+56 more)

### Community 3 - "context_test.go"
Cohesion: 0.16
Nodes (14): RouterContext, FromParams(), NewContext(), setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey(), TestContext_ValidKey(), TestFromParamsTakesOwnershipOfTheLookupSlice() (+6 more)

### Community 4 - "Load tests"
Cohesion: 0.29
Nodes (6): Burst scenarios, Load tests, Running, The vegeta CLI, What each scenario pins, Why a separate module

### Community 6 - "Router"
Cohesion: 0.07
Nodes (30): Chain, Middleware, Router, net/http.Handler, net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer (+22 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 9 - "NewCatalog"
Cohesion: 0.27
Nodes (10): github.com/joaolaureano/go-router/router.Method, Describe(), Catalog(), Route, identity(), NewCatalog(), variablesOf(), wantFor() (+2 more)

### Community 11 - "NewChain"
Cohesion: 0.24
Nodes (8): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestChainSealsOnceItIsBaked(), TestNewChain(), TestNewChain_Middlewares()

### Community 18 - "bench_test.go"
Cohesion: 0.40
Nodes (14): testing.B, BenchmarkCatchAll(), BenchmarkEscaped(), BenchmarkFanout50(), BenchmarkMethodNotAllowed(), BenchmarkMiss(), BenchmarkOptions(), BenchmarkParam() (+6 more)

### Community 20 - "node"
Cohesion: 0.12
Nodes (27): advertisedMethods(), endpoint, lookup, Match, Method, node, cloneEndpoint(), cloneNode() (+19 more)

### Community 28 - "routing/fuzz_test.go"
Cohesion: 0.24
Nodes (14): net/http/httptest.ResponseRecorder, testing.F, assertAllowHeader(), FuzzHeadFollowsGet(), fuzzRouter(), FuzzServeAnyTarget(), Router, serve() (+6 more)

## Knowledge Gaps
- **10 isolated node(s):** `github.com/joaolaureano/go-router/loadtest`, `run.sh script`, `Install`, `Instalação`, `contextKey` (+5 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewRouter()` connect `NewRouter` to `requireNoTransportErrors`, `Router`, `NewCatalog`, `bench_test.go`, `routing/fuzz_test.go`?**
  _High betweenness centrality (0.152) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewCatalog`, `NewRouter`, `NewChain`, `node`?**
  _High betweenness centrality (0.115) - this node is a cross-community bridge._
- **Why does `assertFound()` connect `testing.T` to `node`, `Router`?**
  _High betweenness centrality (0.112) - this node is a cross-community bridge._
- **Are the 56 inferred relationships involving `NewRouter()` (e.g. with `BenchmarkRegister300()` and `benchRouter()`) actually correct?**
  _`NewRouter()` has 56 INFERRED edges - model-reasoned connections that need verification._
- **Are the 44 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `requireNoTransportErrors()` (e.g. with `TestBurst_ColdTree()` and `TestBurst_EveryRouteInTheTable()`) actually correct?**
  _`requireNoTransportErrors()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/joaolaureano/go-router/loadtest`, `run.sh script`, `Install` to the rest of the system?**
  _10 weakly-connected nodes found - possible documentation gaps or missing edges._