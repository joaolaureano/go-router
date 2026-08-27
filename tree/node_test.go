package tree

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNodeInitializesInvariantState(t *testing.T) {
	node := newNode[http.Handler]("users")

	assert.Equal(t, "users", node.path)
	assert.Empty(t, node.children)
	assert.Nil(t, node.parameter)
	assert.Empty(t, node.endpoints)
}

func TestNodeAddChildSeparatesTheThreeSlots(t *testing.T) {
	node := newNode[http.Handler]("users")
	staticChild := newNode[http.Handler]("list")
	parameterChild := newNode[http.Handler](parameterNodePath)
	wildcardChild := newNode[http.Handler](WildcardParam)

	node.addChild(staticChild, staticSegment)
	node.addChild(parameterChild, parameterSegment)
	node.addChild(wildcardChild, wildcardSegment)

	assert.Same(t, staticChild, node.staticChild("list"))
	assert.Nil(t, node.staticChild(parameterNodePath), "a literal segment must not reach the parameter branch")
	assert.Nil(t, node.staticChild(WildcardParam), "a literal segment must not reach the catch-all branch")
	assert.Same(t, parameterChild, node.childFor(parameterNodePath, parameterSegment))
	assert.Same(t, wildcardChild, node.childFor(WildcardParam, wildcardSegment))
	assert.Len(t, node.children, 1)
	assert.Same(t, parameterChild, node.parameter)
	assert.Same(t, wildcardChild, node.wildcard)
}

func TestNodeAddChildRejectsNil(t *testing.T) {
	node := newNode[http.Handler]("users")

	assert.PanicsWithValue(t, "node child must not be nil", func() {
		node.addChild(nil, staticSegment)
	})
}

func TestNodeSetEndpoint(t *testing.T) {
	node := newNode[http.Handler]("users")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	node.setEndpoint(GET, handler, []string{"id"})

	method, exists := node.endpoints[GET]
	assert.True(t, exists)
	assert.NotNil(t, method.handler)
	assert.Equal(t, []string{"id"}, method.variableNames)
}
