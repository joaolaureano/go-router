package tree

import (
	"fmt"
	"net/http"
	"strings"
)

// Tree is the aggregate root. Nothing outside this package holds a node, so
// registration and lookup are the only ways the routing table can change shape.
type Tree struct {
	root *node
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
type Match struct {
	Handler        http.Handler
	Params         []Param
	AllowedMethods []Method
}

func CreateTree() Tree {
	return Tree{
		root: newNode(""),
	}
}

func (t *Tree) RegisterRoute(httpMethod Method, newValue string, method http.Handler) {
	if newValue == "" {
		panic("path must not be empty")
	}
	if method == nil {
		panic("handler must not be nil")
	}
	t.register(httpMethod, newValue, method)
}

func (t *Tree) register(httpMethod Method, path string, method http.Handler) {
	target, variableNames := t.ensurePath(path)
	if target.hasMethod(httpMethod) {
		panic(fmt.Sprintf("Duplicated path: %s", path))
	}
	target.setEndpoint(httpMethod, method, variableNames)
}

// ensurePath walks to the node addressed by path, creating any segment that is
// missing, and reports the variable names the path declares along the way. Those
// names are what an endpoint stored here has to be able to pair with the values
// a lookup collects.
func (t *Tree) ensurePath(path string) (*node, []string) {
	if path[0] != '/' {
		panic("Path must begin with front-slash (/)")
	}
	segments := splitSegments(path)
	if err := validateSegments(path, segments); err != nil {
		panic(err.Error())
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
			nextNode = newNode(nodePath)
			currNode.addChild(nextNode, kind)
		}
		currNode = nextNode
	}
	return currNode, variableNames
}

// Lookup resolves a path against the tree in a single walk.
func (t *Tree) Lookup(httpMethod Method, path string) (Match, Status) {
	search := lookup{httpMethod: httpMethod}
	matched := t.root.match(splitSegments(path), 0, &search)
	if matched == nil {
		if allowed := search.allowedMethods(); allowed != nil {
			return Match{AllowedMethods: allowed}, StatusMethodNotAllowed
		}
		return Match{}, StatusNotFound
	}

	resolved := matched.endpoints[httpMethod]
	return Match{
		Handler: resolved.handler,
		Params:  zipParams(resolved.variableNames, search.values),
	}, StatusFound
}

// zipParams pairs the variable names recorded at registration with the values
// collected on the walk. The walk enters exactly one parameter node per name,
// so the two always line up.
func zipParams(names, values []string) []Param {
	if len(names) == 0 {
		return nil
	}
	params := make([]Param, len(names))
	for i, name := range names {
		params[i] = Param{Name: name, Value: values[i]}
	}
	return params
}

// Merge copies every route of source that this tree does not already define.
func (t *Tree) Merge(source *Tree) {
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
func (t *Tree) MergeAt(prefix string, source *Tree) {
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
			return fmt.Errorf("path contains an empty segment")
		}

		if segment == WildcardParam {
			if i != len(segments)-1 {
				return fmt.Errorf("catch-all '%s' must be the last segment of '%s'", WildcardParam, path)
			}
			continue
		}

		startsWithBrace := segment[0] == '{'
		endsWithBrace := segment[len(segment)-1] == '}'
		if startsWithBrace || endsWithBrace {
			if !startsWithBrace || !endsWithBrace {
				return fmt.Errorf("Delimiter '{' must be closed by '}'")
			}

			paramName := strings.TrimSuffix(strings.TrimPrefix(segment, "{"), "}")
			if paramName == "" || strings.ContainsAny(paramName, "{}") {
				return fmt.Errorf("invalid route parameter: %s", segment)
			}
			if paramName == WildcardParam {
				return fmt.Errorf("'%s' is reserved for the catch-all and cannot name a parameter", WildcardParam)
			}
			if _, exists := paramNames[paramName]; exists {
				return fmt.Errorf("routing pattern '%s' contains duplicate param key, '%s'", path, segment)
			}
			paramNames[paramName] = struct{}{}
		} else if strings.ContainsAny(segment, "{}") {
			return fmt.Errorf("invalid route segment: %s", segment)
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
