package tree

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_const "github.com/joaolaureano/go-router/const"
	"github.com/joaolaureano/go-router/router/context"

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

	assert.PanicsWithValue(t, "path must not be empty", func() {
		tree.RegisterRoute(_const.GET, "", handler)
	})
	assert.Empty(t, tree.root.children, "Invalid registration should not create children")
}

func TestRegister_NilHandler(t *testing.T) {
	tree := CreateTree()

	assert.PanicsWithValue(t, "handler must not be nil", func() {
		tree.RegisterRoute(_const.GET, "/path", nil)
	})
	assert.Empty(t, tree.root.children, "Invalid registration should not create children")
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
	assert.NotNil(t, tree.FindRoute(&context.RouterContext{}, _const.GET, "/"))
	assert.NotNil(t, tree.root.children, "Children should not be nil")
	assert.NotZero(t, len(tree.root.Method), "Method should not be nil")
}

func TestRegister_RootPreservesRoutesAndMethods(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(_const.GET, "/path", handler)
	tree.RegisterRoute(_const.GET, "/", handler)
	tree.RegisterRoute(_const.POST, "/", handler)

	assert.NotNil(t, tree.FindRoute(&context.RouterContext{}, _const.GET, "/path"))
	assert.NotNil(t, tree.FindRoute(&context.RouterContext{}, _const.GET, "/"))
	assert.NotNil(t, tree.FindRoute(&context.RouterContext{}, _const.POST, "/"))
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
	assert.NotNil(t, firstChild.parameter, "Parameter child should not be nil")
	// First child node
	secondChild := firstChild.parameter
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
	secondChild := firstChild.parameter
	assert.NotNil(t, firstChild.parameter, "Parameter child should not be nil")
	assert.NotNil(t, secondChild.children, "First child's children should not be nil")
	assert.Len(t, firstChild.children, 0, "First child should contain static children only")
	assert.Equal(t, "{*}", firstChild.parameter.path, "Path should be {valid}")
	assert.Zero(t, len(firstChild.parameter.Method), "Method should be nil")
	//
	//// Third child node - branching paths
	branch1 := firstChild.parameter.children[0]
	branch2 := firstChild.parameter.children[1]
	assert.NotNil(t, branch1.children, "Second child's children should not be nil")
	assert.NotNil(t, branch2.children, "Second child's children should not be nil")
	assert.Equal(t, "path1", branch1.path, "Path should be /path1")
	assert.Equal(t, "path2", branch2.path, "Path should be /path2")
	assert.NotNil(t, branch1.Method, "Method should not be nil")
	assert.NotNil(t, branch2.Method, "Method should not be nil")
	assert.NotZero(t, len(branch1.Method), "Method should be 1")
	assert.NotZero(t, len(branch2.Method), "Method should be 1")
}

func TestValidatePath_ValidPaths(t *testing.T) {
	validPaths := []string{"valid/{path}", "users/{userID}/posts/{postID}", "health"}
	for _, path := range validPaths {
		assert.NoError(t, validateSegments(path, splitSegments(path)), path)
	}
}

func TestValidatePath_InvalidPaths(t *testing.T) {
	invalidPaths := []string{
		"invalid/{path",
		"users/{}",
		"users/{id}/posts/{id}",
		"users//posts",
		"users/{id}extra",
		"users/{id}/posts}",
	}
	for _, path := range invalidPaths {
		assert.Error(t, validateSegments(path, splitSegments(path)), path)
	}
}

func TestFindRoute(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/path/{test}/test2", handler)
	foundNode := tree.FindRoute(ctx, _const.GET, "/path/test1/test2")
	assert.NotNil(t, foundNode, "Found node should not be nil")
}

func TestFindRoute_Root(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/", handler)
	foundNode := tree.FindRoute(nil, _const.GET, "/")

	assert.NotNil(t, foundNode, "Found node should not be nil")
}

func TestFindRoute_PathEmpty(t *testing.T) {
	tree := CreateTree()
	foundNode := tree.FindRoute(nil, _const.GET, "")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_InexistentRoot(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	foundNode := tree.FindRoute(ctx, _const.GET, "/")
	assert.Nil(t, foundNode, "Found node should be nil")
}

func TestFindRoute_InexistentRootMethod(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	ctx := &context.RouterContext{}
	tree.RegisterRoute(_const.GET, "/", handler)

	foundNode := tree.FindRoute(ctx, _const.POST, "/")
	// Root
	assert.Nil(t, foundNode, "Path should be nil")
}

func TestFindRoute_InexistentPathOnlyRoot(t *testing.T) {
	tree := CreateTree()
	ctx := &context.RouterContext{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(_const.GET, "/", handler)

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

func TestFindRoute_FallsBackToParameterAfterStaticBranchFails(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	ctx := context.NewContext()

	tree.RegisterRoute(_const.GET, "/files/static/view", handler)
	tree.RegisterRoute(_const.GET, "/files/{id}/edit", handler)

	foundNode := tree.FindRoute(ctx, _const.GET, "/files/static/edit")

	assert.NotNil(t, foundNode)
	assert.Equal(t, "static", ctx.Value("id"))
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

func TestTree_Merge(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(_const.GET, "/path/", handler)

	tree.RegisterRoute(_const.GET, "/path/test", handler)

	tree2 := CreateTree()

	tree2.RegisterRoute(_const.GET, "/pathz/", handler)

	tree2.RegisterRoute(_const.GET, "/pathz/test", handler)

	tree.Merge(&tree2)

	assert.NotNil(t, tree.FindRoute(&context.RouterContext{}, _const.GET, "/pathz"))
	assert.NotNil(t, tree.FindRoute(&context.RouterContext{}, _const.GET, "/pathz/test"))
	assert.NotNil(t, tree2.FindRoute(&context.RouterContext{}, _const.GET, "/pathz"))
	assert.NotNil(t, tree2.FindRoute(&context.RouterContext{}, _const.GET, "/pathz/test"))
}

func TestTree_MergeKeepsExistingRoute(t *testing.T) {
	target := CreateTree()
	target.RegisterRoute(_const.GET, "/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("target"))
	}))

	source := CreateTree()
	source.RegisterRoute(_const.GET, "/users/{userID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("source"))
	}))

	target.Merge(&source)
	matched := target.FindRoute(context.NewContext(), _const.GET, "/users/42")
	assert.NotNil(t, matched)

	response := httptest.NewRecorder()
	matched.Method[_const.GET].Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/users/42", nil))
	assert.Equal(t, "target", response.Body.String())
}

func TestTree_MergeCopiesRoutesWithoutRebuildingPaths(t *testing.T) {
	target := CreateTree()
	source := CreateTree()
	source.RegisterRoute(_const.GET, "/users/{id}/posts/{postID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("source"))
	}))
	source.RegisterRoute(_const.POST, "/users/{id}/posts/{postID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("created"))
	}))

	target.Merge(&source)
	ctx := context.NewContext()
	matched := target.FindRoute(ctx, _const.GET, "/users/42/posts/7")

	assert.NotNil(t, matched)
	assert.Equal(t, "42", ctx.Value("id"))
	assert.Equal(t, "7", ctx.Value("postID"))
	assert.NotNil(t, target.FindRoute(context.NewContext(), _const.POST, "/users/42/posts/7"))
}

func TestTree_MergeWithItselfIsNoOp(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(_const.GET, "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	tree.Merge(&tree)

	assert.NotNil(t, tree.FindRoute(context.NewContext(), _const.GET, "/health"))
}
