package tree

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"

	_const "web/const"
	"web/router/context"
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
	currNode := &t.root
	if path[0] != '/' {
		panic("Path must begin with front-slash (/)")
	}
	if path == "/" {
		*currNode = &Node{
			path:     "/",
			children: []*Node{},
			Method:   map[_const.HTTPMethods]Method{httpMethod: {Handler: method, variableName: nil}},
		}
		return
	}
	path = strings.Trim(path, "/")
	validatePath(path)
	pathVariablesName := make([]string, 0)
	for _, pathSplitted := range strings.Split(path, "/") {
		nextNode := (*currNode).getChild(pathSplitted)
		if nextNode == nil {
			nextNode = &Node{
				path:     pathSplitted,
				children: []*Node{},
				Method:   make(map[_const.HTTPMethods]Method),
			}
			if isParam(pathSplitted) {
				nextNode.path = "{*}"
				(*currNode).children = append((*currNode).children, nextNode)
				pathVariablesName = append(pathVariablesName, strings.Trim(pathSplitted, "{}"))
			} else {
				(*currNode).children = append([]*Node{nextNode}, (*currNode).children...)
			}
		}
		currNode = &nextNode
	}
	if !reflect.ValueOf((*(currNode)).Method[httpMethod]).IsZero() {
		panic(fmt.Sprintf("Duplicated path: %s", path))
	}
	(*currNode).setEndpoint(httpMethod, method, pathVariablesName)
}

func (t *Tree) FindRoute(ctx *context.RouterContext, httpMethods _const.HTTPMethods, value string) *Node {
	if len(value) == 0 {
		return nil
	}
	return t.findRoute(ctx, httpMethods, value)
}

func (t *Tree) findRoute(ctx *context.RouterContext, httpMethod _const.HTTPMethods, path string) *Node {
	currNode := t.root
	if len((*currNode).children) == 0 {
		return nil
	}
	paths := strings.Split(strings.Trim(path, "/"), "/")
	idx := 0
	pathVariableValues := make([]string, 0)
	for {
		nextNode := (*currNode).getChild(paths[idx])
		if nextNode == nil {
			return nil
		}
		if isParam(nextNode.path) {
			pathVariableValues = append(pathVariableValues, paths[idx])
		}
		idx++
		if idx == len(paths) {
			if reflect.ValueOf(nextNode.Method[httpMethod]).IsZero() {
				return nil
			}
			setPathVariableValues(ctx, nextNode.Method[httpMethod].variableName, pathVariableValues)
			return nextNode

		}
		currNode = nextNode
	}
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
	idx := slices.IndexFunc(n.children, func(n *Node) bool {
		return path == n.path
	})
	if idx >= 0 {
		return n.children[idx]
	}
	if isParam(n.children[len(n.children)-1].path) {
		return n.children[len(n.children)-1]
	}
	return nil
}

func validatePath(path string) {
	paramList := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
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

func isParam(path string) bool {
	return path[0] == '{' && path[len(path)-1] == '}'
}
