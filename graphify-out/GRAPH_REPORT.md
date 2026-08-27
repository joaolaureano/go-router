# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 23 files · ~14,848 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 324 nodes · 717 edges · 34 communities (11 shown, 23 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 70 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bab4927d`
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
- NewChain
- sync.RWMutex
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- Method
- bench_test.go
- E
- E
- E
- E
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- routing/fuzz_test.go
- Method
- E
- tree.Param
- Method

## God Nodes (most connected - your core abstractions)
1. `CreateTree()` - 48 edges
2. `Router` - 24 edges
3. `node` - 17 edges
4. `assertFound()` - 14 edges
5. `Param()` - 12 edges
6. `FromRequest()` - 12 edges
7. `state` - 12 edges
8. `setup()` - 11 edges
9. `benchRouter()` - 11 edges
10. `run()` - 11 edges

## Surprising Connections (you probably didn't know these)
- `NewPrefixRouter()` --calls--> `CreateTree()`  [EXTRACTED]
  router/router.go → routing/routing.go
- `NewPrefixRouter()` --calls--> `NewPrefixRouter()`  [EXTRACTED]
  gorouter.go → router/router.go
- `TestDecodeSegmentLeavesMalformedEscapingAlone()` --calls--> `decodeSegment()`  [INFERRED]
  router/router_test.go → router/router.go
- `concreteFor()` --calls--> `classify()`  [INFERRED]
  routing/fuzz_test.go → routing/routing.go
- `NewRouter()` --calls--> `NewRouter()`  [EXTRACTED]
  gorouter.go → router/router.go

## Import Cycles
- None detected.

## Communities (34 total, 23 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.09
Nodes (49): CreateTree(), assertFound(), assertPanicsWith(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_CatchAll() (+41 more)

### Community 2 - "testing.T"
Cohesion: 0.06
Nodes (68): net/http/httptest.Server, testing.T, TestNewPrefixRouterRegistersUnderThePrefix(), TestNewRouterSatisfiesTheAcceptSurface(), TestReExportedVocabularyMatchesTheRouter(), Param(), setup(), TestDecodeSegmentLeavesMalformedEscapingAlone() (+60 more)

### Community 3 - "FromRequest"
Cohesion: 0.12
Nodes (23): contextKey, routeContext, RouterContext, context.Context, github.com/joaolaureano/go-router/routing.Param, net/http.Request, FromParams(), FromRequest() (+15 more)

### Community 6 - "Router"
Cohesion: 0.08
Nodes (23): Chain, chain.Middleware, Router, net/http.Handler, net/http.HandlerFunc, net/http.ResponseWriter, sync/atomic.Bool, sync/atomic.Pointer (+15 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "NewChain"
Cohesion: 0.24
Nodes (8): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestChainSealsOnceItIsBaked(), TestNewChain(), TestNewChain_Middlewares()

### Community 18 - "bench_test.go"
Cohesion: 0.44
Nodes (13): testing.B, BenchmarkCatchAll(), BenchmarkEscaped(), BenchmarkFanout50(), BenchmarkMethodNotAllowed(), BenchmarkMiss(), BenchmarkOptions(), BenchmarkParam() (+5 more)

### Community 20 - "E"
Cohesion: 0.11
Nodes (30): E, Method, node, Param, endpoint, lookup, Match, node (+22 more)

### Community 28 - "routing/fuzz_test.go"
Cohesion: 0.24
Nodes (14): net/http/httptest.ResponseRecorder, testing.F, Router, assertAllowHeader(), FuzzHeadFollowsGet(), fuzzRouter(), FuzzServeAnyTarget(), serve() (+6 more)

## Knowledge Gaps
- **4 isolated node(s):** `Install`, `Instalação`, `contextKey`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreateTree()` connect `routing_test.go` to `E`, `Router`?**
  _High betweenness centrality (0.226) - this node is a cross-community bridge._
- **Why does `assertAllowHeader()` connect `routing/fuzz_test.go` to `testing.T`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `NewChain`?**
  _High betweenness centrality (0.104) - this node is a cross-community bridge._
- **Are the 44 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Install`, `Instalação`, `contextKey` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `routing_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08571428571428572 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05875251509054326 - nodes in this community are weakly interconnected._