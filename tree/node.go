package tree

import (
	"net/http"
	"sort"

	_const "github.com/joaolaureano/go-router/const"
)

// node is one path segment of the tree. It is deliberately unexported: the
// invariant tying an endpoint's variable names to the parameter nodes above it
// only holds while Tree is the sole writer.
type node struct {
	path      string
	children  []*node
	parameter *node
	endpoints map[_const.HTTPMethods]endpoint
}

// endpoint is what a single method registered on a node resolves to. The
// variable names live here rather than on the node because two methods on the
// same path may have been written with different parameter names.
type endpoint struct {
	handler       http.Handler
	variableNames []string
}

func newNode(path string) *node {
	return &node{
		path:      path,
		children:  make([]*node, 0),
		endpoints: make(map[_const.HTTPMethods]endpoint),
	}
}

// staticChild finds the child holding this exact segment. It never returns the
// parameter branch: a request segment that happens to read "{*}" is literal
// text, not a request for the parameter slot.
func (n *node) staticChild(path string) *node {
	for _, child := range n.children {
		if path == child.path {
			return child
		}
	}
	return nil
}

// childFor returns the slot a registration should descend into. Only the
// registrar knows whether a segment was written as a parameter, so only it may
// ask for the parameter branch.
func (n *node) childFor(path string, parameter bool) *node {
	if parameter {
		return n.parameter
	}
	return n.staticChild(path)
}

func (n *node) addChild(child *node, parameter bool) {
	if child == nil {
		panic("node child must not be nil")
	}
	if parameter {
		n.parameter = child
		return
	}
	n.children = append(n.children, child)
}

func (n *node) setEndpoint(httpMethod _const.HTTPMethods, handler http.Handler, variableNames []string) {
	n.endpoints[httpMethod] = endpoint{
		handler:       handler,
		variableNames: variableNames,
	}
}

func (n *node) hasMethod(httpMethod _const.HTTPMethods) bool {
	_, exists := n.endpoints[httpMethod]
	return exists
}

func (n *node) hasAnyMethod() bool {
	return len(n.endpoints) > 0
}

// lookup carries the state of one walk. Parameter values accumulate as the
// walk descends and unwind when a branch fails, so only the values on the
// surviving route remain. Nodes that end the path under some other verb are
// recorded on the way, which lets a single walk answer both "which handler"
// and "which methods would have worked".
type lookup struct {
	httpMethod _const.HTTPMethods
	values     []string
	allowed    map[_const.HTTPMethods]struct{}
}

func (search *lookup) recordAllowed(n *node) {
	if !n.hasAnyMethod() {
		return
	}
	if search.allowed == nil {
		search.allowed = make(map[_const.HTTPMethods]struct{}, len(n.endpoints))
	}
	for httpMethod := range n.endpoints {
		search.allowed[httpMethod] = struct{}{}
	}
}

// allowedMethods returns the recorded methods sorted, so that the Allow header
// of a 405 stays stable across responses.
func (search *lookup) allowedMethods() []string {
	if len(search.allowed) == 0 {
		return nil
	}
	methods := make([]string, 0, len(search.allowed))
	for httpMethod := range search.allowed {
		methods = append(methods, string(httpMethod))
	}
	sort.Strings(methods)
	return methods
}

// match walks the remaining segments, preferring the static child and falling
// back to the parameter branch, and returns the first node answering the
// method being searched for.
func (n *node) match(paths []string, index int, search *lookup) *node {
	if index == len(paths) {
		if n.hasMethod(search.httpMethod) {
			return n
		}
		search.recordAllowed(n)
		return nil
	}

	if child := n.staticChild(paths[index]); child != nil {
		if matched := child.match(paths, index+1, search); matched != nil {
			return matched
		}
	}

	if n.parameter != nil {
		search.values = append(search.values, paths[index])
		if matched := n.parameter.match(paths, index+1, search); matched != nil {
			return matched
		}
		search.values = search.values[:len(search.values)-1]
	}

	return nil
}

func mergeNodes(target, source *node) {
	for httpMethod, sourceEndpoint := range source.endpoints {
		if !target.hasMethod(httpMethod) {
			target.endpoints[httpMethod] = cloneEndpoint(sourceEndpoint)
		}
	}

	for _, sourceChild := range source.children {
		targetChild := target.staticChild(sourceChild.path)
		if targetChild == nil {
			target.children = append(target.children, cloneNode(sourceChild))
			continue
		}
		mergeNodes(targetChild, sourceChild)
	}

	if source.parameter == nil {
		return
	}
	if target.parameter == nil {
		target.parameter = cloneNode(source.parameter)
		return
	}
	mergeNodes(target.parameter, source.parameter)
}

func cloneNode(source *node) *node {
	clone := newNode(source.path)
	for httpMethod, sourceEndpoint := range source.endpoints {
		clone.endpoints[httpMethod] = cloneEndpoint(sourceEndpoint)
	}
	for _, child := range source.children {
		clone.children = append(clone.children, cloneNode(child))
	}
	if source.parameter != nil {
		clone.parameter = cloneNode(source.parameter)
	}
	return clone
}

func cloneEndpoint(source endpoint) endpoint {
	return endpoint{
		handler:       source.handler,
		variableNames: append([]string(nil), source.variableNames...),
	}
}
