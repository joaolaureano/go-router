package router

import (
	"net/http"
	"sync"
	"sync/atomic"
)

// state is what a router and everything derived from it -- its groups, its
// Withs -- share: the routing table and the two fallback handlers.
//
// Serving reads it with plain atomic loads and takes no lock at all. Even a
// read lock would not do here: RWMutex.RLock writes to a counter every reader
// on every core touches, and that one cache line caps a router at roughly the
// throughput of a single core no matter how many are serving.
//
// Writers pay for that instead. They serialise on mu and publish by swapping in
// a whole new tree, so a reader is always walking one complete version or the
// next, never one being edited underneath it. Registration is a start-up cost;
// copying the table to change it is the right side of that trade.
type state struct {
	mu               sync.Mutex
	tree             atomic.Pointer[routes]
	notFound         atomic.Pointer[http.HandlerFunc]
	methodNotAllowed atomic.Pointer[http.HandlerFunc]

	// served turns one-way the first time a request arrives. Until then there
	// are no readers to protect, so building the table costs no copies at all;
	// after it, every change is published by copy. Setting it takes mu once,
	// which is what orders it against a writer already inside mutateTree.
	served atomic.Bool
}

func newState(tree *routes) *state {
	notFound := http.HandlerFunc(http.NotFound)
	methodNotAllowed := http.HandlerFunc(defaultMethodNotAllowed)

	shared := &state{}
	shared.tree.Store(tree)
	shared.notFound.Store(&notFound)
	shared.methodNotAllowed.Store(&methodNotAllowed)
	return shared
}

// markServed records that readers now exist. Callers must not hold mu.
func (s *state) markServed() {
	if s.served.Load() {
		return
	}
	s.mu.Lock()
	s.served.Store(true)
	s.mu.Unlock()
}

// mutateTree applies fn to the routing table. Callers must hold mu.
func (s *state) mutateTree(fn func(*routes)) {
	tree := s.tree.Load()
	if !s.served.Load() {
		// Nobody can be walking it yet, so editing in place is safe and spares
		// the copy. Building a table of n routes stays O(n) rather than O(n²).
		fn(tree)
		return
	}
	next := tree.Clone()
	fn(next)
	s.tree.Store(next)
}

func (s *state) setNotFound(handler http.HandlerFunc) {
	s.notFound.Store(&handler)
}

func (s *state) setMethodNotAllowed(handler http.HandlerFunc) {
	s.methodNotAllowed.Store(&handler)
}
