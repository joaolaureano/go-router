package tree

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func assertPanicsWith(t *testing.T, want error, fn func()) {
	t.Helper()
	defer func() {
		recovered := recover()
		err, ok := recovered.(error)
		if assert.True(t, ok, "expected a panic carrying an error, got %v", recovered) {
			assert.ErrorIs(t, err, want)
		}
	}()
	fn()
}

func assertFound(t *testing.T, tree *Tree, httpMethod Method, path string) {
	t.Helper()
	_, status := tree.Lookup(httpMethod, path)
	assert.Equal(t, StatusFound, status, path)
}

func TestCreateTree(t *testing.T) {
	tree := CreateTree()
	assert.NotNil(t, tree.root, "Root should not be nil")
	assert.Empty(t, tree.root.path, "Path should begin empty")
	assert.Empty(t, tree.root.children, "Children should begin empty")
	assert.Zero(t, len(tree.root.endpoints), "Method should be empty")
}

func TestRegister_EmptyPath(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	assert.PanicsWithError(t, ErrEmptyPath.Error(), func() {
		tree.RegisterRoute(GET, "", handler)
	})
	assert.Empty(t, tree.root.children, "Invalid registration should not create children")
}

func TestRegister_NilHandler(t *testing.T) {
	tree := CreateTree()

	assert.PanicsWithError(t, ErrNilHandler.Error(), func() {
		tree.RegisterRoute(GET, "/path", nil)
	})
	assert.Empty(t, tree.root.children, "Invalid registration should not create children")
}

func TestRegister_DuplicatedPathVariable(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	assert.Panics(t, func() { tree.RegisterRoute(GET, "/test/{test}/{test}", handler) }, "Insert should panic for an invalid path")
}

func TestRegister_PanicInvalidPath(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	invalidPath := "invalidPath"

	assertPanicsWith(t, ErrPathNotRooted, func() { tree.RegisterRoute(GET, invalidPath, handler) })
}

func TestRegister_OnlyRoot(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(GET, "/", handler)

	// Root
	assertFound(t, &tree, GET, "/")
	assert.NotNil(t, tree.root.children, "Children should not be nil")
	assert.NotZero(t, len(tree.root.endpoints), "Method should not be nil")
}

func TestRegister_RootPreservesRoutesAndMethods(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(GET, "/path", handler)
	tree.RegisterRoute(GET, "/", handler)
	tree.RegisterRoute(POST, "/", handler)

	assertFound(t, &tree, GET, "/path")
	assertFound(t, &tree, GET, "/")
	assertFound(t, &tree, POST, "/")
}

func TestRegister_SimpleTree(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.RegisterRoute(GET, "/path/{valid}/path", handler)

	// Root
	firstChild := tree.root.children[0]
	assert.NotNil(t, firstChild.children, "Children should not be nil")
	assert.Equal(t, "path", firstChild.path, "Path should be /path")
	assert.Zero(t, len(firstChild.endpoints), "Method should be nil")
	assert.NotNil(t, firstChild.parameter, "Parameter child should not be nil")
	// First child node
	secondChild := firstChild.parameter
	assert.NotNil(t, secondChild.children, "First child's children should not be nil")
	assert.Equal(t, "{*}", secondChild.path, "Path should be {valid}")
	assert.Zero(t, len(secondChild.endpoints), "Method should be nil")
	// Second child node
	thirdChild := secondChild.children[0]
	assert.NotNil(t, thirdChild.children, "Second child's children should not be nil")
	assert.Equal(t, "path", thirdChild.path, "Path should be /path")
	assert.NotZero(t, len(thirdChild.endpoints), "Method should not be empty")
	assert.NotZero(t, thirdChild.endpoints[GET], "Method should not be nil")
}

func TestRegister_DuplicatedPath(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	assert.Panics(t,
		func() {
			tree.RegisterRoute(GET, "/path/{valid}/path", handler)
			tree.RegisterRoute(GET, "/path/{valid}/path", handler)
		}, "Should panic when creating same route multiple times")
}

func TestRegister_MultipleBranches(t *testing.T) {
	tree := CreateTree()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tree.RegisterRoute(GET, "/path/{valid}/path1", handler)
	tree.RegisterRoute(GET, "/path/{valid}/path2", handler)

	// First Child
	firstChild := tree.root.children[0]
	assert.NotNil(t, firstChild.children, "Children should not be nil")
	assert.Equal(t, "path", firstChild.path, "Path should be /path")
	assert.Zero(t, len(firstChild.endpoints), "Method should be nil")

	// Second child
	secondChild := firstChild.parameter
	assert.NotNil(t, firstChild.parameter, "Parameter child should not be nil")
	assert.NotNil(t, secondChild.children, "First child's children should not be nil")
	assert.Len(t, firstChild.children, 0, "First child should contain static children only")
	assert.Equal(t, "{*}", firstChild.parameter.path, "Path should be {valid}")
	assert.Zero(t, len(firstChild.parameter.endpoints), "Method should be nil")
	//
	//// Third child node - branching paths
	branch1 := firstChild.parameter.children[0]
	branch2 := firstChild.parameter.children[1]
	assert.NotNil(t, branch1.children, "Second child's children should not be nil")
	assert.NotNil(t, branch2.children, "Second child's children should not be nil")
	assert.Equal(t, "path1", branch1.path, "Path should be /path1")
	assert.Equal(t, "path2", branch2.path, "Path should be /path2")
	assert.NotNil(t, branch1.endpoints, "Method should not be nil")
	assert.NotNil(t, branch2.endpoints, "Method should not be nil")
	assert.NotZero(t, len(branch1.endpoints), "Method should be 1")
	assert.NotZero(t, len(branch2.endpoints), "Method should be 1")
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

func TestRegister_PanicsCarryDistinguishableErrors(t *testing.T) {
	for name, testCase := range map[string]struct {
		register func(tree *Tree)
		want     error
	}{
		"empty path":          {func(tree *Tree) { tree.RegisterRoute(GET, "", handler) }, ErrEmptyPath},
		"nil handler":         {func(tree *Tree) { tree.RegisterRoute(GET, "/path", nil) }, ErrNilHandler},
		"not rooted":          {func(tree *Tree) { tree.RegisterRoute(GET, "path", handler) }, ErrPathNotRooted},
		"unbalanced brace":    {func(tree *Tree) { tree.RegisterRoute(GET, "/{id", handler) }, ErrInvalidPattern},
		"duplicate param":     {func(tree *Tree) { tree.RegisterRoute(GET, "/{id}/{id}", handler) }, ErrInvalidPattern},
		"misplaced catch-all": {func(tree *Tree) { tree.RegisterRoute(GET, "/*/edit", handler) }, ErrInvalidPattern},
		"duplicate route": {func(tree *Tree) {
			tree.RegisterRoute(GET, "/path", handler)
			tree.RegisterRoute(GET, "/path", handler)
		}, ErrDuplicateRoute},
	} {
		t.Run(name, func(t *testing.T) {
			tree := CreateTree()
			assertPanicsWith(t, testCase.want, func() { testCase.register(&tree) })
		})
	}
}

func TestLookup(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/path", handler)

	match, status := tree.Lookup(GET, "/path")

	assert.Equal(t, StatusFound, status)
	assert.NotNil(t, match.Handler)
	assert.Empty(t, match.Params)
}

func TestLookup_Root(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/", handler)

	for _, path := range []string{"/", ""} {
		match, status := tree.Lookup(GET, path)

		assert.Equal(t, StatusFound, status, path)
		assert.NotNil(t, match.Handler, path)
	}
}

func TestLookup_UnregisteredRoot(t *testing.T) {
	tree := CreateTree()

	_, status := tree.Lookup(GET, "/")

	assert.Equal(t, StatusNotFound, status)
}

func TestLookup_RootRegisteredUnderAnotherMethod(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(POST, "/", handler)

	match, status := tree.Lookup(GET, "/")

	assert.Equal(t, StatusMethodNotAllowed, status)
	assert.Equal(t, []Method{POST}, match.AllowedMethods)
}

func TestLookup_UnknownPath(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/path", handler)

	for _, path := range []string{"/other", "/path/deeper", "/"} {
		match, status := tree.Lookup(GET, path)

		assert.Equal(t, StatusNotFound, status, path)
		assert.Nil(t, match.Handler, path)
		assert.Nil(t, match.AllowedMethods, path)
	}
}

func TestLookup_PathRegisteredUnderAnotherMethod(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(POST, "/path", handler)

	match, status := tree.Lookup(GET, "/path")

	assert.Equal(t, StatusMethodNotAllowed, status)
	assert.Nil(t, match.Handler)
	assert.Equal(t, []Method{POST}, match.AllowedMethods)
}

func TestLookup_IntermediateNodeIsNotAnEndpoint(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/path/leaf", handler)

	_, status := tree.Lookup(GET, "/path")

	assert.Equal(t, StatusNotFound, status)
}

func TestLookup_CapturesParameters(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/users/{userID}/posts/{postID}", handler)

	match, status := tree.Lookup(GET, "/users/7/posts/42")

	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: "userID", Value: "7"}, {Name: "postID", Value: "42"}}, match.Params)
}

func TestLookup_FallsBackToParameterAfterStaticBranchFails(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/static/leaf", handler)
	tree.RegisterRoute(GET, "/{id}/other", handler)

	match, status := tree.Lookup(GET, "/static/other")

	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: "id", Value: "static"}}, match.Params)
}

func TestLookup_DiscardsValuesFromAbandonedBranches(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/{first}/dead-end/leaf", handler)
	tree.RegisterRoute(GET, "/{only}/other", handler)

	match, status := tree.Lookup(GET, "/value/other")

	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: "only", Value: "value"}}, match.Params)
}

func TestLookup_AllowCollectsEveryBranchThatEndsThePath(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(POST, "/a/b", handler)
	tree.RegisterRoute(PATCH, "/a/{id}", handler)

	match, status := tree.Lookup(GET, "/a/b")

	assert.Equal(t, StatusMethodNotAllowed, status)
	assert.Equal(t, []Method{PATCH, POST}, match.AllowedMethods)
}

func TestLookup_CatchAll(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/files/*", handler)

	for path, want := range map[string]string{
		"/files/a/b/c": "a/b/c",
		"/files/a":     "a",
		"/files/":      "",
		"/files":       "",
	} {
		match, status := tree.Lookup(GET, path)

		assert.Equal(t, StatusFound, status, path)
		assert.Equal(t, []Param{{Name: WildcardParam, Value: want}}, match.Params, path)
	}
}

func TestLookup_CatchAllYieldsToStaticAndParameter(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/files/exact", handler)
	tree.RegisterRoute(GET, "/files/{id}/edit", handler)
	tree.RegisterRoute(GET, "/files/*", handler)

	match, status := tree.Lookup(GET, "/files/exact")
	assert.Equal(t, StatusFound, status)
	assert.Empty(t, match.Params, "a static route wins over the catch-all")

	match, status = tree.Lookup(GET, "/files/7/edit")
	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: "id", Value: "7"}}, match.Params, "a parameter route wins over the catch-all")

	match, status = tree.Lookup(GET, "/files/7/other")
	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: WildcardParam, Value: "7/other"}}, match.Params, "the catch-all takes what nothing else matched")
}

func TestLookup_CatchAllAfterParameters(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/{tenant}/files/*", handler)

	match, status := tree.Lookup(GET, "/acme/files/a/b")

	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: "tenant", Value: "acme"}, {Name: WildcardParam, Value: "a/b"}}, match.Params)
}

func TestLookup_CatchAllIsMethodAware(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(POST, "/files/*", handler)

	match, status := tree.Lookup(GET, "/files/a/b")

	assert.Equal(t, StatusMethodNotAllowed, status)
	assert.Equal(t, []Method{POST}, match.AllowedMethods)
}

func TestRegister_CatchAllMustBeLast(t *testing.T) {
	tree := CreateTree()

	assert.Panics(t, func() { tree.RegisterRoute(GET, "/files/*/edit", handler) })
}

func TestRegister_WildcardNameIsReserved(t *testing.T) {
	tree := CreateTree()

	assert.Panics(t, func() { tree.RegisterRoute(GET, "/{*}", handler) })
}

func TestLookup_EmptyTree(t *testing.T) {
	tree := CreateTree()

	for _, path := range []string{"/", "", "/path"} {
		_, status := tree.Lookup(GET, path)

		assert.Equal(t, StatusNotFound, status, path)
	}
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

	tree.RegisterRoute(GET, "/path/", handler)

	tree.RegisterRoute(GET, "/path/test", handler)

	tree2 := CreateTree()

	tree2.RegisterRoute(GET, "/pathz/", handler)

	tree2.RegisterRoute(GET, "/pathz/test", handler)

	tree.Merge(&tree2)

	assertFound(t, &tree, GET, "/pathz")
	assertFound(t, &tree, GET, "/pathz/test")
	assertFound(t, &tree2, GET, "/pathz")
	assertFound(t, &tree2, GET, "/pathz/test")
}

func TestTree_MergeKeepsExistingRoute(t *testing.T) {
	target := CreateTree()
	target.RegisterRoute(GET, "/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("target"))
	}))

	source := CreateTree()
	source.RegisterRoute(GET, "/users/{userID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("source"))
	}))

	target.Merge(&source)
	matched, status := target.Lookup(GET, "/users/42")
	assert.Equal(t, StatusFound, status)

	response := httptest.NewRecorder()
	matched.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/users/42", nil))
	assert.Equal(t, "target", response.Body.String())
}

func TestTree_MergeCopiesRoutesWithoutRebuildingPaths(t *testing.T) {
	target := CreateTree()
	source := CreateTree()
	source.RegisterRoute(GET, "/users/{id}/posts/{postID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("source"))
	}))
	source.RegisterRoute(POST, "/users/{id}/posts/{postID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("created"))
	}))

	target.Merge(&source)
	matched, status := target.Lookup(GET, "/users/42/posts/7")

	assert.Equal(t, StatusFound, status)
	assert.Equal(t, []Param{{Name: "id", Value: "42"}, {Name: "postID", Value: "7"}}, matched.Params)
	assertFound(t, &target, POST, "/users/42/posts/7")
}

func TestTree_MergeWithItselfIsNoOp(t *testing.T) {
	tree := CreateTree()
	tree.RegisterRoute(GET, "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	tree.Merge(&tree)

	assertFound(t, &tree, GET, "/health")
}
