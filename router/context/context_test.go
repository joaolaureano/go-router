package context

import (
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

func setup() *RouterContext {
	return NewContext()
}
