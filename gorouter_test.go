package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joaolaureano/go-router/router"
	"github.com/stretchr/testify/assert"
)

// TestNewRouterSatisfiesTheAcceptSurface pins what this package exists for: a
// caller can hold the narrow interface while the constructor hands back the
// concrete type, and can name a method without importing anything else.
func TestNewRouterSatisfiesTheAcceptSurface(t *testing.T) {
	var r Router = NewRouter()

	r.Register(GET, "/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("pong"))
	})

	response := httptest.NewRecorder()
	r.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ping", nil))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "pong", response.Body.String())
}

func TestNewPrefixRouterRegistersUnderThePrefix(t *testing.T) {
	r := NewPrefixRouter("/api")

	r.Register(GET, "/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	served := httptest.NewRecorder()
	r.ServeHTTP(served, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	assert.Equal(t, "ok", served.Body.String())

	unprefixed := httptest.NewRecorder()
	r.ServeHTTP(unprefixed, httptest.NewRequest(http.MethodGet, "/health", nil))
	assert.Equal(t, http.StatusNotFound, unprefixed.Code, "the prefix is not optional")
}

// TestReExportedVocabularyMatchesTheRouter guards against the constants drifting
// from the ones they stand in for, which no compiler check would catch.
func TestReExportedVocabularyMatchesTheRouter(t *testing.T) {
	assert.Equal(t, router.WildcardParam, WildcardParam)
	assert.Equal(t, []Method{router.GET, router.HEAD, router.POST, router.PUT, router.PATCH, router.DELETE, router.OPTIONS},
		[]Method{GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS})
}
