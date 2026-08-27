package routing

import "errors"

// Registering a bad route is a programming error, caught at start-up rather
// than per request, so it panics rather than returning. The panic value is one
// of these errors, wrapped with the offending pattern, so a caller recovering
// around its route table can tell the cases apart with errors.Is instead of
// matching on message text.
//
// Internal assertions -- states the public API cannot reach -- stay plain
// panics: nobody can act on those.
var (
	ErrEmptyPath      = errors.New("path must not be empty")
	ErrNilHandler     = errors.New("handler must not be nil")
	ErrPathNotRooted  = errors.New("path must begin with a front slash")
	ErrDuplicateRoute = errors.New("route already registered")
	ErrInvalidPattern = errors.New("invalid route pattern")
)
