package routing

import (
	"fmt"
	"slices"
	"strings"
)

// Tree is the aggregate root. Nothing outside this package holds a node, so
// registration and lookup are the only ways the routing table can change shape.
//
// E is whatever a route resolves to. Matching a path against a pattern owes
// nothing to HTTP, so the tree does not name http.Handler; the router pins E
// when it builds one.
type Tree[E any] struct {
	root *node[E]
}

// Status says how a lookup ended. Carrying "the path exists under other
// methods" in the result is what lets one walk answer what used to take two.
type Status int

const (
	StatusNotFound Status = iota
	StatusMethodNotAllowed
	StatusFound
)

// Param is a route variable captured by a lookup.
type Param struct {
	Name  string
	Value string
}

// Match is the outcome of a lookup: a plain value the caller is free to
// interpret. The tree reports what it found and takes no part in deciding how
// that reaches a handler.
type Match[E any] struct {
	Handler        E
	Params         []Param
	AllowedMethods []Method
}

func CreateTree[E any]() Tree[E] {
	return Tree[E]{
		root: newNode[E](""),
	}
}

func (t *Tree[E]) RegisterRoute(httpMethod Method, newValue string, method E) {
	if newValue == "" {
		panic(ErrEmptyPath)
	}
	// E may be an interface, in which case a caller can hand over a nil one.
	// Boxing to any and comparing catches exactly that, and leaves a typed nil
	// alone: that is a handler, just a broken one.
	if any(method) == nil {
		panic(ErrNilHandler)
	}
	t.register(httpMethod, newValue, method)
}

func (t *Tree[E]) register(httpMethod Method, path string, method E) {
	target, variableNames := t.ensurePath(path, false)
	if target.hasMethod(httpMethod) {
		panic(fmt.Errorf("%w: %s %s", ErrDuplicateRoute, httpMethod, path))
	}
	target.setEndpoint(httpMethod, method, variableNames)
}

// ensurePath walks to the node addressed by path, creating any segment that is
// missing, and reports the variable names the path declares along the way. Those
// names are what an endpoint stored here has to be able to pair with the values
// a lookup collects.
//
// Consecutive literal segments are joined with "/" and inserted as one run,
// which is what lets the tree compress a shared prefix across several of them
// into a single edge instead of one node per segment. A parameter always
// interrupts a run: it is inserted separately, so it can only ever anchor to
// the node that run's insertion returns -- a node that represents having fully
// consumed some number of whole segments, never a partial one.
//
// A run immediately followed by a parameter carries that separating "/" as
// its own trailing byte, matching is a lookup arrives with it still attached to
// the raw path: a literal edge is matched by raw byte comparison, with nothing
// else to strip the separator the way a parameter's own capture does for
// whatever follows it. A run with no parameter after it -- the tail of the
// pattern -- needs no such trailing byte, and a parameter with nothing literal
// before it (the first segment, or right after another parameter) needs none
// either, since there is no pending run to close.
//
// tailNeedsSeparator asks for that same trailing byte on the very last run,
// even though nothing in path itself follows it. MergeAt passes true: the node
// this returns becomes the attachment point for a whole grafted subtree, which
// is exactly the position a parameter would be in, and needs the same "/"
// accounted for. A plain registration passes false -- there path really does
// end at the node this returns, with nothing left to separate it from.
func (t *Tree[E]) ensurePath(path string, tailNeedsSeparator bool) (*node[E], []string) {
	if path[0] != '/' {
		panic(fmt.Errorf("%w: %s", ErrPathNotRooted, path))
	}
	segments := splitSegments(path)
	if err := validateSegments(path, segments); err != nil {
		panic(err)
	}

	currNode := t.root
	var variableNames []string
	var literalRun []string
	flushLiteralRun := func(beforeParameter bool) {
		if len(literalRun) == 0 {
			return
		}
		text := strings.Join(literalRun, "/")
		if beforeParameter {
			text += "/"
		}
		currNode = ensureStaticChild(currNode, text)
		literalRun = literalRun[:0]
	}
	for _, segment := range segments {
		if isParam(segment) {
			flushLiteralRun(true)
			variableNames = append(variableNames, strings.Trim(segment, "{}"))
			if currNode.parameter == nil {
				currNode.parameter = newNode[E](parameterNodePath)
			}
			currNode = currNode.parameter
		} else {
			literalRun = append(literalRun, segment)
		}
	}
	flushLiteralRun(tailNeedsSeparator)
	return currNode, variableNames
}

// Lookup resolves a path against the tree in a single walk, allocating nothing
// unless the route captures variables.
func (t *Tree[E]) Lookup(httpMethod Method, path string) (Match[E], Status) {
	return t.LookupDecoded(httpMethod, path, nil)
}

// LookupDecoded is Lookup for a path whose segments still need transforming --
// percent-decoding being the reason it exists. decode runs on each segment as
// the walk reaches it, so nothing has to be materialised up front and a path
// with nothing to decode passes nil and pays nothing.
func (t *Tree[E]) LookupDecoded(httpMethod Method, path string, decode func(string) string) (Match[E], Status) {
	search := lookup{httpMethod: httpMethod}
	trimmed := strings.Trim(path, "/")

	var matched *node[E]
	if decode == nil {
		matched = t.root.match(trimmed, &search)
	} else {
		matched = t.root.matchDecoded(newDecodedCursor(trimmed, decode), &search)
	}
	if matched == nil {
		if allowed := search.allowedMethods(); allowed != nil {
			return Match[E]{AllowedMethods: allowed}, StatusMethodNotAllowed
		}
		return Match[E]{}, StatusNotFound
	}

	resolved := matched.find(httpMethod)
	return Match[E]{
		Handler: resolved.handler,
		Params:  nameParams(search.params, resolved.variableNames),
	}, StatusFound
}

// nameParams labels the values the walk captured with the names the matched
// endpoint recorded at registration. The walk enters exactly one parameter node
// per name, so the two always line up.
func nameParams(params []Param, names []string) []Param {
	if len(names) == 0 {
		return nil
	}
	for i, name := range names {
		params[i].Name = name
	}
	return params
}

// Clone returns a deep copy sharing only the endpoints themselves. It is what
// lets a caller publish changes without locking readers: mutate the copy, then
// swap it in, and a reader is always walking one whole tree or the other.
func (t *Tree[E]) Clone() *Tree[E] {
	return &Tree[E]{root: cloneNode(t.root, nil)}
}

// Merge copies every route of source that this tree does not already define.
func (t *Tree[E]) Merge(source *Tree[E]) {
	if t.root == source.root {
		return
	}
	mergeNodes(t.root, source.root, nil)
}

// MergeAt copies source's routes in under prefix.
//
// The prefix may itself declare variables. A grafted endpoint recorded its own
// names against its own depth, so those of the prefix are prepended: without
// that, a lookup would pair the prefix's captured value with the first name the
// grafted route declared.
func (t *Tree[E]) MergeAt(prefix string, source *Tree[E]) {
	target, prefixVariables := t.ensurePath(prefix, true)
	if target == source.root {
		return
	}
	mergeNodes(target, source.root, prefixVariables)
}

// splitSegments breaks a path into the segments the tree is keyed by. Only
// registration uses it: a lookup walks the path with a cursor instead, since
// materialising the segments is the one allocation a static route would
// otherwise pay. Both sides have to agree on this normalisation, or a route
// gets stored under a shape no request will ever reach.
//
// The root carries no segments, so it needs no special case on either side:
// an empty segment list simply leaves the walk standing on the root node.
func splitSegments(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func validateSegments(path string, segments []string) error {
	// A slice, not a set: a path declares a handful of variables, and scanning a
	// few short strings costs less than the map this used to allocate on every
	// registration.
	var paramNames []string
	for _, segment := range segments {
		if segment == "" {
			return fmt.Errorf("%w: %s contains an empty segment", ErrInvalidPattern, path)
		}

		startsWithBrace := segment[0] == '{'
		endsWithBrace := segment[len(segment)-1] == '}'
		if startsWithBrace || endsWithBrace {
			if !startsWithBrace || !endsWithBrace {
				return fmt.Errorf("%w: unbalanced braces in segment %q of %s", ErrInvalidPattern, segment, path)
			}

			paramName := strings.TrimSuffix(strings.TrimPrefix(segment, "{"), "}")
			if paramName == "" || strings.ContainsAny(paramName, "{}") {
				return fmt.Errorf("%w: invalid parameter %q in %s", ErrInvalidPattern, segment, path)
			}
			if slices.Contains(paramNames, paramName) {
				return fmt.Errorf("%w: %s declares %q twice", ErrInvalidPattern, path, segment)
			}
			paramNames = append(paramNames, paramName)
		} else if strings.ContainsAny(segment, "{}") {
			return fmt.Errorf("%w: invalid segment %q in %s", ErrInvalidPattern, segment, path)
		}
	}
	return nil
}

// parameterNodePath is the internal path every parameter node is stored under.
// It is not a pattern anyone writes: braces anywhere but wrapping a whole
// segment are rejected, so "{*}" as a request segment is always literal text.
const parameterNodePath = "{*}"

func isParam(path string) bool {
	return len(path) >= 3 && path[0] == '{' && path[len(path)-1] == '}'
}
