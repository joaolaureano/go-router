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
	node.children = append([]*Node{child}, node.children...)
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

func matchPath(node *Node, paths []string, index int, values []string) (*Node, []string) {
	if index == len(paths) {
		if len(node.Method) == 0 {
			return nil, nil
		}
		return node, values
	}

	if child := node.getChild(paths[index]); child != nil {
		if matchedNode, matchedValues := matchPath(child, paths, index+1, values); matchedNode != nil {
			return matchedNode, matchedValues
		}
	}

	if node.parameter != nil {
		values = append(values, paths[index])
		if matchedNode, matchedValues := matchPath(node.parameter, paths, index+1, values); matchedNode != nil {
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
