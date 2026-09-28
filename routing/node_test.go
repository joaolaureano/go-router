package routing

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

func TestEnsureStaticChildCreatesALeafOnAnEmptyNode(t *testing.T) {
	node := newNode[http.Handler]("")

	child := ensureStaticChild(node, "users")

	assert.Len(t, node.children, 1)
	assert.Same(t, child, node.children[0])
	assert.Equal(t, "users", child.path)
}

func TestEnsureStaticChildSplitsOnTheFirstDivergingByte(t *testing.T) {
	node := newNode[http.Handler]("")

	first := ensureStaticChild(node, "path1")
	second := ensureStaticChild(node, "path2")

	// "path1" and "path2" share "path": a patricia trie holds that once, as one
	// intermediate edge, rather than two full, separately-stored strings.
	assert.Len(t, node.children, 1)
	branchPoint := node.children[0]
	assert.Equal(t, "path", branchPoint.path)
	assert.Len(t, branchPoint.children, 2)
	assert.Same(t, first, branchPoint.children[0])
	assert.Same(t, second, branchPoint.children[1])
	assert.Equal(t, "1", first.path)
	assert.Equal(t, "2", second.path)
}

func TestEnsureStaticChildReturnsTheSameNodeForTheSameKey(t *testing.T) {
	node := newNode[http.Handler]("")

	first := ensureStaticChild(node, "users")
	second := ensureStaticChild(node, "users")

	assert.Same(t, first, second)
}

func TestEnsureStaticChildDescendsPastAFullyConsumedEdge(t *testing.T) {
	node := newNode[http.Handler]("")

	ensureStaticChild(node, "user")
	agent := ensureStaticChild(node, "userAgent")

	assert.Len(t, node.children, 1, "\"user\" is a prefix of \"userAgent\": no split is needed, only a deeper child")
	assert.Equal(t, "user", node.children[0].path)
	assert.Same(t, agent, node.children[0].children[0])
	assert.Equal(t, "Agent", agent.path)
}

func TestNodeParameterIsASeparateSlotFromStaticChildren(t *testing.T) {
	node := newNode[http.Handler]("users")
	staticChild := ensureStaticChild(node, "list")
	node.parameter = newNode[http.Handler](parameterNodePath)

	_, foundByFirstByte := node.childByFirstByte('l')
	assert.Same(t, staticChild, foundByFirstByte)
	assert.NotSame(t, node.parameter, staticChild)
}

func TestNodeSetEndpoint(t *testing.T) {
	node := newNode[http.Handler]("users")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	node.setEndpoint(GET, handler, []string{"id"})

	registered := node.find(GET)
	assert.NotNil(t, registered)
	assert.NotNil(t, registered.handler)
	assert.Equal(t, []string{"id"}, registered.variableNames)
	assert.Nil(t, node.find(POST), "an unregistered method finds nothing")
}

func TestNodeSetEndpointOnlyAppends(t *testing.T) {
	node := newNode[http.Handler]("users")

	node.setEndpoint(GET, handler, []string{"id"})
	node.setEndpoint(POST, handler, nil)

	assert.Len(t, node.endpoints, 2, "a method already recorded never reaches here: register rejects a duplicate and merge skips one")
	assert.Equal(t, []string{"id"}, node.find(GET).variableNames)
	assert.Nil(t, node.find(POST).variableNames)
}
