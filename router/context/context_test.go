package context

import (
	"context"
	"net/http/httptest"
	"testing"

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

func TestInjectIntoRequest(t *testing.T) {
	routerCtx := &RouterContext{
		// Initialize RouterContext properties for testing if needed
	}

	req := httptest.NewRequest("GET", "/", nil)
	routerCtx.InjectIntoRequest(req)

	// Retrieve the RouterContext from the request's context
	ctxValue := req.Context().Value(RouterContextKey)
	injectedCtx, ok := ctxValue.(*RouterContext)
	if !ok {
		t.Errorf("Expected RouterContext, got %T", ctxValue)
	}
	assert.NotNil(t, injectedCtx)
	if injectedCtx == nil {
		t.Error("Injected context is nil")
	}
	if injectedCtx != routerCtx {
		t.Error("Injected context does not match the original context")
	}

}

func TestWithRequestDoesNotMutateOriginalRequest(t *testing.T) {
	routerCtx := NewContext()
	request := httptest.NewRequest("GET", "/", nil)
	derivedRequest := routerCtx.WithRequest(request)

	assert.Nil(t, request.Context().Value(RouterContextKey))
	assert.Same(t, routerCtx, derivedRequest.Context().Value(RouterContextKey))
}

func TestFromRequest(t *testing.T) {
	routerCtx := NewContext()
	request := routerCtx.WithRequest(httptest.NewRequest("GET", "/", nil))

	foundContext, ok := FromRequest(request)

	assert.True(t, ok)
	assert.Same(t, routerCtx, foundContext)
}

func TestFromRequestSupportsLegacyContextKey(t *testing.T) {
	routerCtx := NewContext()
	request := httptest.NewRequest("GET", "/", nil).WithContext(
		context.WithValue(context.Background(), RouterContextKey, routerCtx),
	)

	foundContext, ok := FromRequest(request)

	assert.True(t, ok)
	assert.Same(t, routerCtx, foundContext)
}

func setup() *RouterContext {
	return NewContext()
}
