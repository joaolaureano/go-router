# Graph Report - go-router  (2026-08-27)

## Corpus Check
- 18 files · ~10,003 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 244 nodes · 584 edges · 35 communities (12 shown, 23 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 102 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d9ca8fab`
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
- E
- github.com/joaolaureano/go-router/router/context.RouterContext
- github.com/joaolaureano/go-router/tree.RouterTree
- Go-Router
- github.com/joaolaureano/go-router/router.Router
- Method
- NewChain
- node
- E
- Method
- node_test.go
- Method
- github.com/joaolaureano/go-router/tree.Method
- Method
- github.com/joaolaureano/go-router/tree.Param
- Method
- node
- segmentKind
- E
- tree.Param
- Method

## God Nodes (most connected - your core abstractions)
1. `NewRouter()` - 52 edges
2. `CreateTree()` - 37 edges
3. `Router` - 26 edges
4. `node` - 19 edges
5. `Method` - 16 edges
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
  router/router.go → routing/routing.go
- `advertisedMethods()` --references--> `Method`  [EXTRACTED]
  router/router.go → routing/method.go

## Import Cycles
- None detected.

## Communities (35 total, 23 thin omitted)

### Community 0 - "routing_test.go"
Cohesion: 0.11
Nodes (37): CreateTree(), assertFound(), assertPanicsWith(), TestCreateTree(), TestLookup(), TestLookup_AllowCollectsEveryBranchThatEndsThePath(), TestLookup_CapturesParameters(), TestLookup_CatchAll() (+29 more)

### Community 2 - "testing.T"
Cohesion: 0.11
Nodes (54): net/http/httptest.Server, testing.T, Param(), NewRouter(), setup(), TestNewRouter(), TestRouter_AnswersOptionsAutomatically(), TestRouter_BacktracksWhenStaticMatchLacksMethod() (+46 more)

### Community 3 - "FromRequest"
Cohesion: 0.21
Nodes (14): contextKey, RouterContext, FromParams(), FromRequest(), NewContext(), setup(), TestContext_InvalidKey(), TestContext_SetOverwritesExistingKey() (+6 more)

### Community 6 - "Router"
Cohesion: 0.18
Nodes (10): net/http.HandlerFunc, net/http.Request, net/http.ResponseWriter, sync.RWMutex, Router, advertisedMethods(), defaultMethodNotAllowed(), requestSegments() (+2 more)

### Community 7 - "Repository"
Cohesion: 0.39
Nodes (3): Article, Repository, New()

### Community 11 - "Chain"
Cohesion: 0.17
Nodes (9): Chain, Middleware, Router, net/http.Handler, NewPrefixRouter(), NewRouter(), NewPrefixRouter(), TestNewRouterWithPrefix() (+1 more)

### Community 17 - "Method"
Cohesion: 0.21
Nodes (14): Match, Method, classify(), E, isParam(), splitSegments(), TestIsParam(), TestValidatePath_InvalidPaths() (+6 more)

### Community 18 - "NewChain"
Cohesion: 0.24
Nodes (7): NewChain(), TestChain_MiddlewareRunsInRegistrationOrder(), TestChain_MultipleMiddleware(), TestChain_NoMiddleware(), TestChain_SingleMiddleware(), TestNewChain(), TestNewChain_Middlewares()

### Community 20 - "node"
Cohesion: 0.29
Nodes (11): endpoint, lookup, node, cloneEndpoint(), cloneNode(), node[E], E, mergeBranch() (+3 more)

### Community 23 - "node_test.go"
Cohesion: 0.40
Nodes (4): TestNewNodeInitializesInvariantState(), TestNodeAddChildRejectsNil(), TestNodeAddChildSeparatesTheThreeSlots(), TestNodeSetEndpoint()

## Knowledge Gaps
- **4 isolated node(s):** `contextKey`, `Install`, `Instalação`, `github.com/joaolaureano/go-router`
  These have ≤1 connection - possible missing edges or undocumented components.
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `assertFound()` connect `routing_test.go` to `Method`, `testing.T`, `Chain`?**
  _High betweenness centrality (0.127) - this node is a cross-community bridge._
- **Why does `Router` connect `Router` to `testing.T`, `NewChain`, `Chain`?**
  _High betweenness centrality (0.119) - this node is a cross-community bridge._
- **Why does `Method` connect `Method` to `routing_test.go`, `node`, `Router`?**
  _High betweenness centrality (0.118) - this node is a cross-community bridge._
- **Are the 49 inferred relationships involving `NewRouter()` (e.g. with `TestNewRouter()` and `TestRouter_AnswersOptionsAutomatically()`) actually correct?**
  _`NewRouter()` has 49 INFERRED edges - model-reasoned connections that need verification._
- **Are the 33 inferred relationships involving `CreateTree()` (e.g. with `TestCreateTree()` and `TestLookup()`) actually correct?**
  _`CreateTree()` has 33 INFERRED edges - model-reasoned connections that need verification._
- **What connects `contextKey`, `Install`, `Instalação` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `routing_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10810810810810811 - nodes in this community are weakly interconnected._