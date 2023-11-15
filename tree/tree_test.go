package tree

import (
	"net/http"
	"testing"

	_const "web/const"
	"web/router/context"

	"github.com/stretchr/testify/assert"
)

func TestCreateTree(t *testing.T) {
	tree := CreateTree()
	assert.NotNil(t, tree.root, "Root should not be nil")
	assert.Empty(t, tree.root.path, "Path should begin empty")
	assert.Empty(t, tree.root.children, "Children should begin empty")
	assert.Zero(t, len(tree.root.Method), "Method should be empty")
}

func TestRegister_EmptyPath(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(_const.GET, "", handler)

	// Root
	assert.Len(t, tree.root.children, 0, "Children should not be nil")
}

func TestRegister_DuplicatedPathVariable(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	assert.Panics(t, func() { tree.RegisterRoute(_const.GET, "/test/{test}/{test}", handler) }, "Insert should panic for an invalid path")
}

func TestRegister_PanicInvalidPath(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	invalidPath := "invalidPath"

	assert.PanicsWithValue(t, "Path must begin with front-slash (/)", func() { tree.RegisterRoute(_const.GET, invalidPath, handler) }, "Insert should panic for an invalid path")
}

func TestRegister_OnlyRoot(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(_const.GET, "/", handler)

	// Root
	assert.Equal(t, "/", tree.root.path, "Path should be /")
	assert.NotNil(t, tree.root.children, "Children should not be nil")
	assert.NotZero(t, len(tree.root.Method), "Method should not be nil")
}

func TestRegister_SimpleTree(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(_const.GET, "/path/{valid}/path", handler)

	// Root
	firstChild := tree.root.children[0]
	assert.NotNil(t, firstChild.children, "Children should not be nil")
	assert.Equal(t, "path", firstChild.path, "Path should be /path")
	assert.Zero(t, len(firstChild.Method), "Method should be nil")
	// First child node
	secondChild := firstChild.children[0]
	assert.NotNil(t, secondChild.children, "First child's children should not be nil")
	assert.Equal(t, "{*}", secondChild.path, "Path should be {valid}")
	assert.Zero(t, len(secondChild.Method), "Method should be nil")
	// Second child node
	thirdChild := secondChild.children[0]
	assert.NotNil(t, thirdChild.children, "Second child's children should not be nil")
	assert.Equal(t, "path", thirdChild.path, "Path should be /path")
	assert.NotZero(t, len(thirdChild.Method), "Method should not be empty")
	assert.NotZero(t, thirdChild.Method[_const.GET], "Method should not be nil")
}

func TestRegister_DuplicatedPath(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	assert.Panics(t,
		func() {
			tree.RegisterRoute(_const.GET, "/path/{valid}/path", handler)
			tree.RegisterRoute(_const.GET, "/path/{valid}/path", handler)
		}, "Should panic when creating same route multiple times")
}

func TestRegister_MultipleBranches(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path/{valid}/path1", handler)
	tree.RegisterRoute(_const.GET, "/path/{valid}/path2", handler)

	// First Child
	firstChild := tree.root.children[0]
	assert.NotNil(t, firstChild.children, "Children should not be nil")
	assert.Equal(t, "path", firstChild.path, "Path should be /path")
	assert.Zero(t, len(firstChild.Method), "Method should be nil")

	// Second child
	secondChild := firstChild.children[0]
	assert.NotNil(t, secondChild.children, "First child's children should not be nil")
	assert.Len(t, secondChild.children, 2, "First child's children should not be nil")
	assert.Equal(t, "{*}", secondChild.path, "Path should be {valid}")
	assert.Zero(t, len(secondChild.Method), "Method should be nil")
	//
	//// Third child node - branching paths
	branch1 := secondChild.children[1]
	branch2 := secondChild.children[0]
	assert.NotNil(t, branch1.children, "Second child's children should not be nil")
	assert.NotNil(t, branch2.children, "Second child's children should not be nil")
	assert.Equal(t, "path1", branch1.path, "Path should be /path1")
	assert.Equal(t, "path2", branch2.path, "Path should be /path2")
	assert.NotNil(t, branch1.Method, "Method should not be nil")
	assert.NotNil(t, branch2.Method, "Method should not be nil")
	assert.NotZero(t, len(branch1.Method), "Method should be 1")
	assert.NotZero(t, len(branch2.Method), "Method should be 1")
}

func TestValidatePath_ShouldNotPanic(t *testing.T) {
	assert.NotPanics(t, func() { validatePath("valid/{path}") }, "Should not panics with valid path")
}

func TestValidatePath_ShouldPanic(t *testing.T) {
	assert.Panics(t, func() { validatePath("invalid/{path") }, "Should panic with invalid path")
}

func TestFindRoute(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path/{test}/test2", handler)
	foundNode := tree.FindRoute(ctx, _const.GET, "/path/test1/test2")
	assert.NotNil(t, foundNode, "Found node should not be nil")
}

func TestFindRoute_PathEmpty(t *testing.T) {
	tree := CreateTree()
	foundNode := tree.FindRoute(nil, _const.GET, "")
	assert.Nil(t, foundNode, "Found node should not be nil")
}

func TestFindRoute_InexistentRoot(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	foundNode := tree.FindRoute(ctx, _const.GET, "/path/error")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_InexistentPath(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path", handler)
	foundNode := tree.FindRoute(ctx, _const.GET, "/path/error")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_InexistentMethod(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path", handler)
	foundNode := tree.FindRoute(ctx, _const.POST, "/path")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_RootWithoutChildren(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	foundNode := tree.FindRoute(ctx, _const.GET, "")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_EmptyPath(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path", handler)
	foundNode := tree.FindRoute(ctx, _const.GET, "")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_InvalidPath(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path", handler)
	foundNode := tree.FindRoute(ctx, _const.GET, "/")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestIsParam(t *testing.T) {
	assert.True(t, isParam("{param}"), "Should return true for a parameter")
	assert.False(t, isParam("{test"), "Should return false for a non-parameter")
	assert.False(t, isParam("test}"), "Should return false for a non-parameter")
	assert.False(t, isParam("test"), "Should return false for a non-parameter")
}
