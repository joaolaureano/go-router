package tree

import (
	"net/http"
	"sort"

	_const "github.com/joaolaureano/go-router/const"
)

type Node struct {
	path      string
	children  []*Node
	parameter *Node
	Method    map[_const.HTTPMethods]Method
}

type Method struct {
	Handler      http.Handler
	variableName []string
}

func newNode(path string) *Node {
	return &Node{
		path:     path,
		children: make([]*Node, 0),
		Method:   make(map[_const.HTTPMethods]Method),
	}
}

// staticChild finds the child holding this exact segment. It never returns the
// parameter branch: a request segment that happens to read "{*}" is literal
// text, not a request for the parameter slot.
func (node *Node) staticChild(path string) *Node {
	for _, child := range node.children {
		if path == child.path {
			return child
		}
	}
	return nil
}

// childFor returns the slot a registration should descend into. Only the
// registrar knows whether a segment was written as a parameter, so only it may
// ask for the parameter branch.
func (node *Node) childFor(path string, parameter bool) *Node {
	if parameter {
		return node.parameter
	}
	return node.staticChild(path)
}

func (node *Node) addChild(child *Node, parameter bool) {
	if child == nil {
		panic("node child must not be nil")
	}
	if parameter {
		node.parameter = child
		return
	}
	node.children = append(node.children, child)
}

func (node *Node) setEndpoint(httpMethod _const.HTTPMethods, handler http.Handler, pathVariables []string) {
	node.Method[httpMethod] = Method{
		Handler:      handler,
		variableName: pathVariables,
	}
}

func (node *Node) hasMethod(httpMethod _const.HTTPMethods) bool {
	_, exists := node.Method[httpMethod]
	return exists
}

func (node *Node) hasAnyMethod() bool {
	return len(node.Method) > 0
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

func (search *lookup) recordAllowed(node *Node) {
	if !node.hasAnyMethod() {
		return
	}
	if search.allowed == nil {
		search.allowed = make(map[_const.HTTPMethods]struct{}, len(node.Method))
	}
	for httpMethod := range node.Method {
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
func (node *Node) match(paths []string, index int, search *lookup) *Node {
	if index == len(paths) {
		if node.hasMethod(search.httpMethod) {
			return node
		}
		search.recordAllowed(node)
		return nil
	}

	if child := node.staticChild(paths[index]); child != nil {
		if matched := child.match(paths, index+1, search); matched != nil {
			return matched
		}
	}

	if node.parameter != nil {
		search.values = append(search.values, paths[index])
		if matched := node.parameter.match(paths, index+1, search); matched != nil {
			return matched
		}
		search.values = search.values[:len(search.values)-1]
	}

	return nil
}

func mergeNodes(target, source *Node) {
	for httpMethod, method := range source.Method {
		if !target.hasMethod(httpMethod) {
			target.Method[httpMethod] = cloneMethod(method)
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

func cloneNode(source *Node) *Node {
	clone := newNode(source.path)
	for httpMethod, method := range source.Method {
		clone.Method[httpMethod] = cloneMethod(method)
	}
	for _, child := range source.children {
		clone.children = append(clone.children, cloneNode(child))
	}
	if source.parameter != nil {
		clone.parameter = cloneNode(source.parameter)
	}
	return clone
}

func cloneMethod(method Method) Method {
	return Method{
		Handler:      method.Handler,
		variableName: append([]string(nil), method.variableName...),
	}
}
