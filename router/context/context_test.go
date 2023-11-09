package context

import (
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

func setup() *RouterContext {
	return NewContext()
}
