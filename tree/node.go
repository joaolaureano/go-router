package tree

import (
	"net/http"
	"sort"
)

// node is one path segment of the tree. It is deliberately unexported: the
// invariant tying an endpoint's variable names to the parameter nodes above it
// only holds while Tree is the sole writer.
type node struct {
	path      string
	children  []*node
	parameter *node
	endpoints map[Method]endpoint
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
		endpoints: make(map[Method]endpoint),
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

func (n *node) setEndpoint(httpMethod Method, handler http.Handler, variableNames []string) {
	n.endpoints[httpMethod] = endpoint{
		handler:       handler,
		variableNames: variableNames,
	}
}

func (n *node) hasMethod(httpMethod Method) bool {
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
	httpMethod Method
	values     []string
	allowed    map[Method]struct{}
}

func (search *lookup) recordAllowed(n *node) {
	if !n.hasAnyMethod() {
		return
	}
	if search.allowed == nil {
		search.allowed = make(map[Method]struct{}, len(n.endpoints))
	}
	for httpMethod := range n.endpoints {
		search.allowed[httpMethod] = struct{}{}
	}
}

// allowedMethods returns the recorded methods sorted, so that a caller
// rendering them -- into an Allow header, say -- gets a stable order.
func (search *lookup) allowedMethods() []Method {
	if len(search.allowed) == 0 {
		return nil
	}
	methods := make([]Method, 0, len(search.allowed))
	for httpMethod := range search.allowed {
		methods = append(methods, httpMethod)
	}
	sort.Slice(methods, func(i, j int) bool { return methods[i] < methods[j] })
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

// mergeNodes copies source into target. prefixVariables names the variables
// that target already sits below, which every grafted endpoint has to inherit.
func mergeNodes(target, source *node, prefixVariables []string) {
	for httpMethod, sourceEndpoint := range source.endpoints {
		if !target.hasMethod(httpMethod) {
			target.endpoints[httpMethod] = cloneEndpoint(sourceEndpoint, prefixVariables)
		}
	}

	for _, sourceChild := range source.children {
		targetChild := target.staticChild(sourceChild.path)
		if targetChild == nil {
			target.children = append(target.children, cloneNode(sourceChild, prefixVariables))
			continue
		}
		mergeNodes(targetChild, sourceChild, prefixVariables)
	}

	if source.parameter == nil {
		return
	}
	if target.parameter == nil {
		target.parameter = cloneNode(source.parameter, prefixVariables)
		return
	}
	mergeNodes(target.parameter, source.parameter, prefixVariables)
}

func cloneNode(source *node, prefixVariables []string) *node {
	clone := newNode(source.path)
	for httpMethod, sourceEndpoint := range source.endpoints {
		clone.endpoints[httpMethod] = cloneEndpoint(sourceEndpoint, prefixVariables)
	}
	for _, child := range source.children {
		clone.children = append(clone.children, cloneNode(child, prefixVariables))
	}
	if source.parameter != nil {
		clone.parameter = cloneNode(source.parameter, prefixVariables)
	}
	return clone
}

func cloneEndpoint(source endpoint, prefixVariables []string) endpoint {
	variableNames := make([]string, 0, len(prefixVariables)+len(source.variableNames))
	variableNames = append(variableNames, prefixVariables...)
	variableNames = append(variableNames, source.variableNames...)
	return endpoint{
		handler:       source.handler,
		variableNames: variableNames,
	}
}
