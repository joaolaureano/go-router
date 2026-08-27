package router

import (
	"errors"

	"github.com/joaolaureano/go-router/tree"
)

// Misconfiguring a router is a programming error, caught while the route table
// is built rather than per request, so it panics with one of these rather than
// returning. Recover around the setup and use errors.Is to tell them apart.
var (
	// ErrNilHandler is the tree's, not a second error meaning the same thing:
	// a caller recovering around registration sees one value either way.
	ErrNilHandler = tree.ErrNilHandler

	ErrChainSealed       = errors.New("middleware cannot be added once the chain has built a route")
	ErrNilRouter         = errors.New("router must not be nil")
	ErrSharedRoutingTree = errors.New("router must not be mounted onto one that shares its routing tree")
)
