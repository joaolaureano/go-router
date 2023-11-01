package router

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"web/router/context"
)

type Node struct {
	path     string
	children []*Node
	method   map[HTTPMethods]Method
}

type Method struct {
	handler      http.Handler
	variableName []string
}

type Tree struct {
	root *Node
}

type TreeRouter interface {
	RegisterRoute(httpMethod HTTPMethods, newValue string, method http.Handler)
	FindRoute(ctx *context.Context, httpMethods HTTPMethods, value string) *Node
}

func createTree() Tree {
	return Tree{
		root: &Node{
			path:     "",
			children: make([]*Node, 0),
			method:   make(map[HTTPMethods]Method),
		},
	}
}

func (t *Tree) RegisterRoute(httpMethod HTTPMethods, newValue string, method http.Handler) {
	if len(newValue) > 0 {
		t.register(httpMethod, newValue, method)
	}
}

func (t *Tree) register(httpMethod HTTPMethods, path string, method http.Handler) {
	currNode := &t.root
	if path[0] != '/' {
		panic("Path must begin with front-slash (/)")
	}
	if path == "/" {
		if reflect.ValueOf((*currNode).method[httpMethod]).IsZero() {
			*currNode = &Node{
				path:     "/",
				children: []*Node{},
				method:   map[HTTPMethods]Method{httpMethod: {handler: method, variableName: nil}},
			}
		} else {
			(*currNode).method[httpMethod] = Method{handler: method, variableName: nil}
		}

		return
	}
	path = strings.Trim(path, "/")
	validatePath(path)
	pathVariablesName := make([]string, 0)
	for _, pathSplitted := range strings.Split(path, "/") {
		nextNode := (*currNode).getChild(pathSplitted)
		if nextNode == nil {
			//methodMap :=
			nextNode = &Node{
				path:     pathSplitted,
				children: []*Node{},
				method:   make(map[HTTPMethods]Method),
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
	if !reflect.ValueOf((*(currNode)).method[httpMethod]).IsZero() {
		panic(fmt.Sprintf("Duplicated path: %s", path))
	}
	(*currNode).setEndpoint(httpMethod, method, pathVariablesName)
}

func (t *Tree) FindRoute(ctx *context.Context, httpMethods HTTPMethods, value string) *Node {
	if len(value) == 0 {
		return nil
	}
	return t.findRoute(ctx, httpMethods, value)
}

func (t *Tree) findRoute(ctx *context.Context, httpMethod HTTPMethods, path string) *Node {
	currNode := t.root
	if len((*currNode).children) == 0 {
		return nil
	}
	paths := strings.Split(strings.Trim(path, "/"), "/")
	if len(paths) == 0 {
		return nil
	}
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
			if reflect.ValueOf(nextNode.method[httpMethod]).IsZero() {
				return nil
			}
			setPathVariableValues(ctx, nextNode.method[httpMethod].variableName, pathVariableValues)
			return nextNode

		}
		currNode = nextNode
	}
}

func setPathVariableValues(ctx *context.Context, keys, values []string) {
	for i, k := range keys {
		(*ctx).Set(k, values[i])
	}
}

func (n *Node) setEndpoint(httpMethod HTTPMethods, handler http.Handler, pathVariables []string) {
	n.method[httpMethod] = Method{
		handler:      handler,
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
