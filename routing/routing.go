package routing

import (
	"fmt"
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
	target, variableNames := t.ensurePath(path)
	if target.hasMethod(httpMethod) {
		panic(fmt.Errorf("%w: %s %s", ErrDuplicateRoute, httpMethod, path))
	}
	target.setEndpoint(httpMethod, method, variableNames)
}

// ensurePath walks to the node addressed by path, creating any segment that is
// missing, and reports the variable names the path declares along the way. Those
// names are what an endpoint stored here has to be able to pair with the values
// a lookup collects.
func (t *Tree[E]) ensurePath(path string) (*node[E], []string) {
	if path[0] != '/' {
		panic(fmt.Errorf("%w: %s", ErrPathNotRooted, path))
	}
	segments := splitSegments(path)
	if err := validateSegments(path, segments); err != nil {
		panic(err)
	}

	currNode := t.root
	var variableNames []string
	for _, segment := range segments {
		kind := classify(segment)
		nodePath := segment
		switch kind {
		case parameterSegment:
			nodePath = parameterNodePath
			variableNames = append(variableNames, strings.Trim(segment, "{}"))
		case wildcardSegment:
			variableNames = append(variableNames, WildcardParam)
		}
		nextNode := currNode.childFor(nodePath, kind)
		if nextNode == nil {
			nextNode = newNode[E](nodePath)
			currNode.addChild(nextNode, kind)
		}
		currNode = nextNode
	}
	return currNode, variableNames
}

// Lookup resolves a path against the tree in a single walk.
func (t *Tree[E]) Lookup(httpMethod Method, path string) (Match[E], Status) {
	return t.LookupSegments(httpMethod, splitSegments(path))
}

// LookupSegments is Lookup for a caller that has already split the path, which
// is what one has to do to transform segments first -- decoding percent escapes
// being the reason that exists.
func (t *Tree[E]) LookupSegments(httpMethod Method, segments []string) (Match[E], Status) {
	search := lookup{httpMethod: httpMethod}
	matched := t.root.match(segments, 0, &search)
	if matched == nil {
		if allowed := search.allowedMethods(); allowed != nil {
			return Match[E]{AllowedMethods: allowed}, StatusMethodNotAllowed
		}
		return Match[E]{}, StatusNotFound
	}

	resolved := matched.endpoints[httpMethod]
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
	target, prefixVariables := t.ensurePath(prefix)
	if target == source.root {
		return
	}
	mergeNodes(target, source.root, prefixVariables)
}

// splitSegments breaks a path into the segments the tree is keyed by.
// Registration and lookup have to agree on this normalisation, otherwise a
// route can be stored under a shape that no request will ever reach.
//
// The root carries no segments, so it needs no special case on either side:
// an empty segment list simply leaves the walk standing on the root node.
// SplitPath breaks a path into the segments the tree is keyed by. A caller that
// must transform segments before matching them splits with this and then calls
// LookupSegments, so that both sides stay on one normalisation.
func SplitPath(path string) []string {
	return splitSegments(path)
}

func splitSegments(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func validateSegments(path string, segments []string) error {
	paramNames := make(map[string]struct{})
	for i, segment := range segments {
		if segment == "" {
			return fmt.Errorf("%w: %s contains an empty segment", ErrInvalidPattern, path)
		}

		if segment == WildcardParam {
			if i != len(segments)-1 {
				return fmt.Errorf("%w: catch-all %q must be the last segment of %s", ErrInvalidPattern, WildcardParam, path)
			}
			continue
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
			if paramName == WildcardParam {
				return fmt.Errorf("%w: %q is reserved for the catch-all and cannot name a parameter", ErrInvalidPattern, WildcardParam)
			}
			if _, exists := paramNames[paramName]; exists {
				return fmt.Errorf("%w: %s declares %q twice", ErrInvalidPattern, path, segment)
			}
			paramNames[paramName] = struct{}{}
		} else if strings.ContainsAny(segment, "{}") {
			return fmt.Errorf("%w: invalid segment %q in %s", ErrInvalidPattern, segment, path)
		}
	}
	return nil
}

// WildcardParam is the name a catch-all segment captures under. A request to
// "/files/a/b" against "/files/*" reads "a/b" from it.
const WildcardParam = "*"

// parameterNodePath is the internal path every parameter node is stored under.
// It is not a pattern anyone writes: "{*}" as a route parameter is rejected,
// and as a request segment it is literal text.
const parameterNodePath = "{*}"

// classify reports how a pattern segment was written.
func classify(segment string) segmentKind {
	switch {
	case segment == WildcardParam:
		return wildcardSegment
	case isParam(segment):
		return parameterSegment
	default:
		return staticSegment
	}
}

func isParam(path string) bool {
	return len(path) >= 3 && path[0] == '{' && path[len(path)-1] == '}'
}
