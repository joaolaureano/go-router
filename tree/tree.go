package tree

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	_const "github.com/joaolaureano/go-router/const"
	"github.com/joaolaureano/go-router/router/context"
)

type Node struct {
	path     string
	children []*Node
	Method   map[_const.HTTPMethods]Method
}

type Method struct {
	Handler      http.Handler
	variableName []string
}

type Tree struct {
	root *Node
}

type RouterTree interface {
	RegisterRoute(httpMethod _const.HTTPMethods, newValue string, method http.Handler)
	FindRoute(ctx *context.RouterContext, httpMethods _const.HTTPMethods, value string) *Node
	FindPath(value string) *Node
	Merge(tree RouterTree)
	Root() *Node
}

func CreateTree() Tree {
	return Tree{
		root: &Node{
			path:     "",
			children: make([]*Node, 0),
			Method:   make(map[_const.HTTPMethods]Method),
		},
	}
}

func (t *Tree) RegisterRoute(httpMethod _const.HTTPMethods, newValue string, method http.Handler) {
	if len(newValue) > 0 {
		t.register(httpMethod, newValue, method)
	}
}

func (t *Tree) register(httpMethod _const.HTTPMethods, path string, method http.Handler) {
	currNode := t.root
	if path[0] != '/' {
		panic("Path must begin with front-slash (/)")
	}
	if path == "/" {
		if _, exists := currNode.Method[httpMethod]; exists {
			panic(fmt.Sprintf("Duplicated path: %s", path))
		}
		currNode.path = "/"
		currNode.Method[httpMethod] = Method{Handler: method}
		return
	}
	path = strings.Trim(path, "/")
	validatePath(path)
	pathVariablesName := make([]string, 0)
	for _, pathSplitted := range strings.Split(path, "/") {
		nextNode := currNode.getChild(pathSplitted)
		if nextNode == nil {
			nextNode = &Node{
				path:     pathSplitted,
				children: []*Node{},
				Method:   make(map[_const.HTTPMethods]Method),
			}
			if isParam(pathSplitted) {
				currNode.children = append(currNode.children, nextNode)
			} else {
				currNode.children = append([]*Node{nextNode}, currNode.children...)
			}
		}
		if isParam(pathSplitted) {
			nextNode.path = "{*}"
			pathVariablesName = append(pathVariablesName, strings.Trim(pathSplitted, "{}"))
		}
		currNode = nextNode
	}
	if _, exists := currNode.Method[httpMethod]; exists {
		panic(fmt.Sprintf("Duplicated path: %s", path))
	}
	currNode.setEndpoint(httpMethod, method, pathVariablesName)
}

func (t *Tree) FindRoute(ctx *context.RouterContext, httpMethods _const.HTTPMethods, value string) *Node {
	if len(value) == 0 {
		return nil
	}
	node, values := t.findPath(value)
	if node == nil {
		return nil
	}
	routeMethod, exists := node.Method[httpMethods]
	if !exists {
		return nil
	}
	if ctx != nil {
		setPathVariableValues(ctx, routeMethod.variableName, values)
	}
	return node
}

func (t *Tree) FindPath(path string) *Node {
	node, _ := t.findPath(path)
	return node
}

func (t *Tree) findPath(path string) (*Node, []string) {
	currNode := t.root
	if path == "/" || path == "" {
		if currNode.path != "/" && len(currNode.Method) == 0 {
			return nil, nil
		}
		return currNode, nil
	}
	if len(currNode.children) == 0 {
		return nil, nil
	}
	paths := strings.Split(strings.Trim(path, "/"), "/")
	idx := 0
	pathVariableValues := make([]string, 0, len(paths))
	nextNode := currNode.getChild(paths[idx])
	for {
		if nextNode == nil {
			return nil, nil
		}
		if isParam(nextNode.path) {
			pathVariableValues = append(pathVariableValues, paths[idx])
		}
		idx++
		if idx == len(paths) {
			return nextNode, pathVariableValues
		}
		currNode = nextNode
		nextNode = currNode.getChild(paths[idx])
	}
}

func (t *Tree) Merge(tree RouterTree) {
	type entry struct {
		node *Node
		path string
	}
	root := tree.Root()
	stack := []entry{{node: root, path: "/"}}
	for len(stack) > 0 {
		currentEntry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		current := currentEntry.node
		for httpMethod, method := range current.Method {
			t.RegisterRoute(httpMethod, namedPath(currentEntry.path, method.variableName), method.Handler)
		}
		for _, child := range current.children {
			childPath := currentEntry.path + "/" + child.path
			if currentEntry.path == "/" {
				childPath = "/" + child.path
			}
			childPath = strings.TrimRight(childPath, "/")
			stack = append(stack, entry{node: child, path: childPath})
		}
	}
}

func namedPath(path string, variableNames []string) string {
	if path == "/" {
		return path
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	variableIndex := 0
	for i, part := range parts {
		if part == "{*}" && variableIndex < len(variableNames) {
			parts[i] = "{" + variableNames[variableIndex] + "}"
			variableIndex++
		}
	}
	return "/" + strings.Join(parts, "/")
}

func (t *Tree) Root() *Node {
	return t.root
}

func setPathVariableValues(ctx *context.RouterContext, keys, values []string) {
	for i, k := range keys {
		(*ctx).Set(k, values[i])
	}
}

func (n *Node) setEndpoint(httpMethod _const.HTTPMethods, handler http.Handler, pathVariables []string) {
	n.Method[httpMethod] = Method{
		Handler:      handler,
		variableName: pathVariables,
	}
}

func (n *Node) getChild(path string) *Node {

	if len(n.children) == 0 {
		return nil
	}
	for _, child := range n.children {
		if path == child.path {
			return child
		}
	}
	if isParam(n.children[len(n.children)-1].path) {
		return n.children[len(n.children)-1]
	}
	return nil
}

func validatePath(path string) {
	paramList := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if len(v) >= 3 {

			if (v[0] == '{') != (v[len(v)-1] == '}') {
				panic("Delimiter '{' must be closed by '}'")
			}
			if isParam(v) {
				if slices.Contains(paramList, v) {
					panic(fmt.Sprintf("routing pattern '%s' contains duplicate param key, '%s'", path, v))
				}
				paramList = append(paramList, v)
			}
		}
	}
}

func isParam(path string) bool {
	return len(path) >= 3 && path[0] == '{' && path[len(path)-1] == '}'
}
