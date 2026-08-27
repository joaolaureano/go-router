package tree

import (
	"fmt"
	"net/http"
	"strings"

	_const "github.com/joaolaureano/go-router/const"
	"github.com/joaolaureano/go-router/router/context"
)

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
		root: newNode(""),
	}
}

func (t *Tree) RegisterRoute(httpMethod _const.HTTPMethods, newValue string, method http.Handler) {
	if newValue == "" {
		panic("path must not be empty")
	}
	if method == nil {
		panic("handler must not be nil")
	}
	t.register(httpMethod, newValue, method)
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
		currNode.Method[httpMethod] = Method{Handler: method}
		return
	}
	path = strings.Trim(path, "/")
	if err := validatePath(path); err != nil {
		panic(err.Error())
	}
	var pathVariablesName []string
	for _, pathSplitted := range strings.Split(path, "/") {
		isParameter := isParam(pathSplitted)
		nodePath := pathSplitted
		if isParameter {
			nodePath = "{*}"
			pathVariablesName = append(pathVariablesName, strings.Trim(pathSplitted, "{}"))
		}
		nextNode := currNode.getChild(nodePath)
		if nextNode == nil {
			nextNode = newNode(nodePath)
			currNode.addChild(nextNode, isParameter)
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
	if path == "/" || path == "" {
		if len(t.root.Method) == 0 {
			return nil, nil
		}
		return t.root, nil
	}
	if len(t.root.children) == 0 && t.root.parameter == nil {
		return nil, nil
	}
	paths := strings.Split(strings.Trim(path, "/"), "/")
	return matchPath(t.root, paths, 0, nil)
}

func (t *Tree) Merge(tree RouterTree) {
	sourceRoot := tree.Root()
	if t.root == sourceRoot {
		return
	}
	mergeNodes(t.root, sourceRoot)
}

func (t *Tree) Root() *Node {
	return t.root
}

func setPathVariableValues(ctx *context.RouterContext, keys, values []string) {
	for i, k := range keys {
		(*ctx).Set(k, values[i])
	}
}

func validatePath(path string) error {
	paramNames := make(map[string]struct{})
	for _, segment := range strings.Split(path, "/") {
		if segment == "" {
			return fmt.Errorf("path contains an empty segment")
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

func isParam(path string) bool {
	return len(path) >= 3 && path[0] == '{' && path[len(path)-1] == '}'
}
