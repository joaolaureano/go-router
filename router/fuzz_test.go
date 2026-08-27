package router

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fuzzMethods are drawn from by index rather than fuzzed as free text, so that
// the fuzzer spends its budget on paths instead of rediscovering that "GET" is
// spelled GET -- and so that HEAD and OPTIONS, the two verbs answered on a
// route's behalf, come up as often as the rest.
var fuzzMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
}

func fuzzRouter() *Router {
	r := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	r.Get("/", handler)
	r.Get("/users", handler)
	r.Get("/users/{id}", handler)
	r.Post("/users/{id}", handler)
	r.Delete("/users/{id}", handler)
	r.Get("/users/{id}/posts/{postID}", handler)
	r.Head("/head-only", handler)
	r.Options("/explicit-options", handler)
	r.Get("/files/*", handler)
	return r
}

// serve runs one request, reporting false when the target is not something a
// client could have sent in the first place.
func serve(r *Router, method, target string) (*httptest.ResponseRecorder, bool) {
	if !strings.HasPrefix(target, "/") {
		return nil, false
	}
	request, err := http.NewRequest(method, "http://example.com"+target, nil)
	if err != nil {
		return nil, false
	}
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	return response, true
}

// FuzzServeAnyTarget asserts that every request a client can send is answered
// with one of the four outcomes this router has, and that the Allow header it
// sets alongside a refusal is well formed.
func FuzzServeAnyTarget(f *testing.F) {
	for _, seed := range []string{
		"/", "/users", "/users/1", "/users/1/posts/2", "/files/a/b",
		"/%7B*%7D", "/{*}", "/users/a%2Fb", "/users/%zz", "//users//1",
		"/head-only", "/explicit-options", "/users/1/", strings.Repeat("/a", 64),
	} {
		for verb := range len(fuzzMethods) {
			f.Add(uint8(verb), seed)
		}
	}

	r := fuzzRouter()

	f.Fuzz(func(t *testing.T, verb uint8, target string) {
		method := fuzzMethods[int(verb)%len(fuzzMethods)]
		response, sent := serve(r, method, target)
		if !sent {
			t.Skip("not a target a client could send")
		}

		switch response.Code {
		case http.StatusOK, http.StatusNotFound:
		case http.StatusNoContent:
			require.Equal(t, http.MethodOptions, method, "204 is only ever the automatic OPTIONS answer")
			// OPTIONS belongs in its own Allow: the router does answer it, which
			// is exactly what this response is.
			allowed := assertAllowHeader(t, response, method, target)
			assert.Contains(t, allowed, http.MethodOptions, "the answer to OPTIONS has to advertise OPTIONS")
		case http.StatusMethodNotAllowed:
			allowed := assertAllowHeader(t, response, method, target)
			assert.NotContains(t, allowed, method, "Allow %v names the method that was just refused", allowed)
		default:
			t.Fatalf("%s %q answered %d", method, target, response.Code)
		}
	})
}

// assertAllowHeader checks what RFC 9110 asks of Allow whenever it is set: a
// non-empty list, in a stable order, naming each method once. What may appear
// in it depends on why it was set, so that is left to the caller.
func assertAllowHeader(t *testing.T, response *httptest.ResponseRecorder, method, target string) []string {
	t.Helper()

	header := response.Header().Get("Allow")
	require.NotEmpty(t, header, "%s %q refused without saying what would work", method, target)

	advertised := strings.Split(header, ", ")
	assert.True(t, slices.IsSorted(advertised), "Allow %q is not in a stable order", header)
	for i, allowed := range advertised {
		assert.NotEmpty(t, allowed, "Allow %q has an empty entry", header)
		assert.NotContains(t, advertised[i+1:], allowed, "Allow %q names %s twice", header, allowed)
	}
	return advertised
}

// FuzzHeadFollowsGet asserts the fallback RFC 9110 asks for: HEAD is GET
// without a body, so wherever GET reaches a route HEAD has to reach one too,
// and the two must agree on the outcome.
func FuzzHeadFollowsGet(f *testing.F) {
	for _, seed := range []string{
		"/", "/users", "/users/1", "/files/a", "/head-only",
		"/explicit-options", "/missing", "/users/a%2Fb",
	} {
		f.Add(seed)
	}

	r := fuzzRouter()

	f.Fuzz(func(t *testing.T, target string) {
		getResponse, sent := serve(r, http.MethodGet, target)
		if !sent {
			t.Skip("not a target a client could send")
		}
		headResponse, _ := serve(r, http.MethodHead, target)

		if getResponse.Code == http.StatusOK {
			assert.Equal(t, http.StatusOK, headResponse.Code,
				"GET %q was served, so HEAD has to be too", target)
		}
		if getResponse.Code == http.StatusNotFound {
			assert.Equal(t, http.StatusNotFound, headResponse.Code,
				"GET %q found nothing, so neither can HEAD", target)
		}
	})
}
