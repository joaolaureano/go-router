package context

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/joaolaureano/go-router/routing"

	"github.com/stretchr/testify/assert"
)

func TestContext_ValidKey(t *testing.T) {
	ctx := setup()
	ctx.Set("test1", "1")
	ctx.Set("test2", "2")

	value1 := ctx.Value("test1")
	value2 := ctx.Value("test2")

	assert.Equal(t, "1", value1)
	assert.Equal(t, "2", value2)
}

func TestContext_InvalidKey(t *testing.T) {
	ctx := setup()

	value := ctx.Value("test")

	assert.Equal(t, "", value)
}

func TestContext_SetOverwritesExistingKey(t *testing.T) {
	ctx := setup()
	ctx.Set("id", "first")
	ctx.Set("id", "second")

	assert.Equal(t, "second", ctx.Value("id"))
}

func TestWithRequestDoesNotMutateOriginalRequest(t *testing.T) {
	routerCtx := NewContext()
	request := httptest.NewRequest("GET", "/", nil)
	derivedRequest := routerCtx.WithRequest(request)

	_, attachedToOriginal := FromRequest(request)
	assert.False(t, attachedToOriginal)

	derivedContext, ok := FromRequest(derivedRequest)
	assert.True(t, ok)
	assert.Same(t, routerCtx, derivedContext)
}

func TestFromRequest(t *testing.T) {
	routerCtx := NewContext()
	request := routerCtx.WithRequest(httptest.NewRequest("GET", "/", nil))

	foundContext, ok := FromRequest(request)

	assert.True(t, ok)
	assert.Same(t, routerCtx, foundContext)
}

func TestFromRequestWithoutRoutingContext(t *testing.T) {
	foundContext, ok := FromRequest(httptest.NewRequest("GET", "/", nil))

	assert.False(t, ok)
	assert.Nil(t, foundContext)
}

func TestWithParams(t *testing.T) {
	request := WithParams(httptest.NewRequest("GET", "/", nil), []routing.Param{
		{Name: "id", Value: "42"},
	})

	routerCtx, ok := FromRequest(request)

	assert.True(t, ok)
	assert.Equal(t, "42", routerCtx.Value("id"))
	assert.Equal(t, "42", Param(request, "id"))
}

func TestWithParamsDelegatesTheRestOfTheContext(t *testing.T) {
	type otherKey struct{}
	parent := context.WithValue(context.Background(), otherKey{}, "kept")
	request := WithParams(httptest.NewRequest("GET", "/", nil).WithContext(parent), nil)

	assert.Equal(t, "kept", request.Context().Value(otherKey{}), "an unrelated value must still resolve")
	assert.Nil(t, request.Context().Done(), "the embedded context answers the rest of the interface")
}

func TestParam(t *testing.T) {
	routerCtx := NewContext()
	routerCtx.Set("id", "42")
	request := routerCtx.WithRequest(httptest.NewRequest("GET", "/", nil))

	assert.Equal(t, "42", Param(request, "id"))
	assert.Equal(t, "", Param(request, "missing"))
	assert.Equal(t, "", Param(httptest.NewRequest("GET", "/", nil), "id"))
}

func setup() *RouterContext {
	return NewContext()
}

func TestFromParamsTakesOwnershipOfTheLookupSlice(t *testing.T) {
	params := []routing.Param{{Name: "id", Value: "7"}, {Name: "postID", Value: "9"}}

	routerCtx := FromParams(params)

	assert.Equal(t, "7", routerCtx.Value("id"))
	assert.Equal(t, "9", routerCtx.Value("postID"))
	assert.Equal(t, "", routerCtx.Value("absent"))

	routerCtx.Set("id", "8")
	assert.Equal(t, "8", params[0].Value, "the slice is taken over, not copied")
}

func TestValueToleratesANilContext(t *testing.T) {
	// The usual `routerCtx, _ := FromRequest(r)` shape leaves a nil pointer on a
	// route that captured nothing, and reading from it has to stay safe.
	var absent *RouterContext

	assert.Equal(t, "", absent.Value("id"))
}
