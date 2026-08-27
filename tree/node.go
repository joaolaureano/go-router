package tree

import (
	"net/http"

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

func (node *Node) getChild(path string) *Node {
	for _, child := range node.children {
		if path == child.path {
			return child
		}
	}
	if path == "{*}" {
		return node.parameter
	}
	return nil
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

// methodFilter decides which nodes are allowed to terminate a search. Matching
// must be method-aware: a node that only answers POST cannot end a GET lookup,
// otherwise the search stops there instead of backtracking into a sibling
// parameter branch that would have matched.
type methodFilter struct {
	httpMethod _const.HTTPMethods
	anyMethod  bool
}

func filterByMethod(httpMethod _const.HTTPMethods) methodFilter {
	return methodFilter{httpMethod: httpMethod}
}

func filterByAnyMethod() methodFilter {
	return methodFilter{anyMethod: true}
}

func (filter methodFilter) accepts(node *Node) bool {
	if filter.anyMethod {
		return node.hasAnyMethod()
	}
	return node.hasMethod(filter.httpMethod)
}

func matchPath(node *Node, paths []string, index int, values []string, filter methodFilter) (*Node, []string) {
	if index == len(paths) {
		if !filter.accepts(node) {
			return nil, nil
		}
		return node, values
	}

	if child := node.getChild(paths[index]); child != nil {
		if matchedNode, matchedValues := matchPath(child, paths, index+1, values, filter); matchedNode != nil {
			return matchedNode, matchedValues
		}
	}

	if node.parameter != nil {
		values = append(values, paths[index])
		if matchedNode, matchedValues := matchPath(node.parameter, paths, index+1, values, filter); matchedNode != nil {
			return matchedNode, matchedValues
		}
	}

	return nil, nil
}

func mergeNodes(target, source *Node) {
	for httpMethod, method := range source.Method {
		if !target.hasMethod(httpMethod) {
			target.Method[httpMethod] = cloneMethod(method)
		}
	}

	for _, sourceChild := range source.children {
		targetChild := staticChild(target, sourceChild.path)
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

func staticChild(node *Node, path string) *Node {
	for _, child := range node.children {
		if child.path == path {
			return child
		}
	}
	return nil
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
