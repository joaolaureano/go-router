package tree

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNodeInitializesInvariantState(t *testing.T) {
	node := newNode("users")

	assert.Equal(t, "users", node.path)
	assert.Empty(t, node.children)
	assert.Nil(t, node.parameter)
	assert.Empty(t, node.endpoints)
}

func TestNodeAddChildSeparatesStaticAndParameterChildren(t *testing.T) {
	node := newNode("users")
	staticChild := newNode("list")
	parameterChild := newNode("{*}")

	node.addChild(staticChild, false)
	node.addChild(parameterChild, true)

	assert.Same(t, staticChild, node.staticChild("list"))
	assert.Nil(t, node.staticChild("{*}"), "a literal segment must not reach the parameter branch")
	assert.Same(t, parameterChild, node.childFor("{*}", true))
	assert.Len(t, node.children, 1)
	assert.Same(t, parameterChild, node.parameter)
}

func TestNodeAddChildRejectsNil(t *testing.T) {
	node := newNode("users")

	assert.PanicsWithValue(t, "node child must not be nil", func() {
		node.addChild(nil, false)
	})
}

func TestNodeSetEndpoint(t *testing.T) {
	node := newNode("users")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	node.setEndpoint(GET, handler, []string{"id"})

	method, exists := node.endpoints[GET]
	assert.True(t, exists)
	assert.NotNil(t, method.handler)
	assert.Equal(t, []string{"id"}, method.variableNames)
}
