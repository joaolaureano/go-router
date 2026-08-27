package tree

import (
	"net/http"
	"testing"

	_const "github.com/joaolaureano/go-router/const"
	"github.com/stretchr/testify/assert"
)

func TestNewNodeInitializesInvariantState(t *testing.T) {
	node := newNode("users")

	assert.Equal(t, "users", node.path)
	assert.Empty(t, node.children)
	assert.Nil(t, node.parameter)
	assert.Empty(t, node.Method)
}

func TestNodeAddChildSeparatesStaticAndParameterChildren(t *testing.T) {
	node := newNode("users")
	staticChild := newNode("list")
	parameterChild := newNode("{*}")

	node.addChild(staticChild, false)
	node.addChild(parameterChild, true)

	assert.Same(t, staticChild, node.getChild("list"))
	assert.Same(t, parameterChild, node.getChild("{*}"))
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

	node.setEndpoint(_const.GET, handler, []string{"id"})

	method, exists := node.Method[_const.GET]
	assert.True(t, exists)
	assert.NotNil(t, method.Handler)
	assert.Equal(t, []string{"id"}, method.variableName)
}
