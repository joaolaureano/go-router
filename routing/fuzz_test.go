package routing

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// capturedValue is what every variable in a fuzzed pattern is given. Any single
// segment would do; what matters is that the value is the same everywhere, so
// that a test failure is about which name a value was paired with rather than
// about the value itself.
const capturedValue = "value"

// tryRegister reports whether the tree accepted this pattern. Rejecting an
// invalid one by panicking is the contract, so a panic here is not a finding --
// the fuzzer is only interested in what happens to the patterns that get in.
func tryRegister(tree *Tree[string], pattern string) (accepted bool) {
	defer func() {
		if recover() != nil {
			accepted = false
		}
	}()
	tree.RegisterRoute(GET, pattern, "handler")
	return true
}

// concreteFor turns a pattern into a request path that must match it, and lists
// the variables it declares in the order a lookup should report them. It is the
// substitution a client performs when it builds a URL from a documented route.
func concreteFor(pattern string) (path string, names []string) {
	trimmed := strings.Trim(pattern, "/")
	if trimmed == "" {
		return "/", nil
	}
	segments := strings.Split(trimmed, "/")
	for i, segment := range segments {
		switch classify(segment) {
		case parameterSegment:
			names = append(names, strings.Trim(segment, "{}"))
			segments[i] = capturedValue
		case wildcardSegment:
			names = append(names, WildcardParam)
			segments[i] = capturedValue
		}
	}
	return "/" + strings.Join(segments, "/"), names
}

// FuzzRegisteredPatternMatchesItsOwnPath asserts the round trip every routing
// table rests on: substitute a value for each variable of a pattern the tree
// accepted, and a request for the result has to arrive at that pattern, with
// the values paired to the names the pattern declared, in order.
//
// It is the invariant that a literal "{*}" request segment once broke by
// reaching the parameter slot, and that a mounted route once broke by pairing
// its first name with the prefix's value.
func FuzzRegisteredPatternMatchesItsOwnPath(f *testing.F) {
	for _, seed := range []string{
		"/",
		"/users",
		"/users/{id}",
		"/users/{id}/posts/{postID}",
		"/files/*",
		"/*",
		"/{*}",
		"/a%2Fb",
		"/{id}/{id}",
		"/users/{id}/",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, pattern string) {
		tree := CreateTree[string]()
		if !tryRegister(&tree, pattern) {
			t.Skip("pattern rejected at registration")
		}

		path, names := concreteFor(pattern)
		match, status := tree.Lookup(GET, path)

		require.Equal(t, StatusFound, status, "%q registered but %q did not reach it", pattern, path)
		require.Len(t, match.Params, len(names), "%q declares %d variables", pattern, len(names))
		for i, name := range names {
			assert.Equal(t, name, match.Params[i].Name, "variable %d of %q", i, pattern)
			assert.Equal(t, capturedValue, match.Params[i].Value, "variable %q of %q", name, pattern)
		}
	})
}

// FuzzLookupOfAnyPathTerminates asserts that an arbitrary request path is
// answered rather than crashing the walk, whatever backtracking it drives
// through the static, parameter and catch-all branches of a table that has all
// three at every level.
func FuzzLookupOfAnyPathTerminates(f *testing.F) {
	for _, seed := range []string{
		"/", "", "//", "/a/b/c", "/{*}", "/*", "/x/{*}/y",
		strings.Repeat("/a", 64), strings.Repeat("/", 64), "/a%2Fb",
	} {
		f.Add(seed)
	}

	tree := CreateTree[string]()
	for _, pattern := range []string{
		"/", "/a", "/a/b", "/a/{id}", "/a/{id}/c", "/{first}", "/{first}/{second}", "/{first}/*", "/*",
	} {
		tree.RegisterRoute(GET, pattern, "handler")
	}
	tree.RegisterRoute(POST, "/a/b", "handler")

	f.Fuzz(func(t *testing.T, path string) {
		match, status := tree.Lookup(GET, path)

		switch status {
		case StatusFound:
			assert.Equal(t, "handler", match.Handler, "a found route must carry its handler")
			assert.Nil(t, match.AllowedMethods, "a found route allows nothing extra")
		case StatusMethodNotAllowed:
			assert.NotEmpty(t, match.AllowedMethods, "405 has to name the methods that would work")
		case StatusNotFound:
			assert.Empty(t, match.Params, "a miss captures nothing")
		default:
			t.Fatalf("unknown status %d for %q", status, path)
		}
	})
}
