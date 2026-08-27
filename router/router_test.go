package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/joaolaureano/go-router/router/context"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, router.root, "Root should not be nil")

}
func TestNewRouterWithPrefix(t *testing.T) {
	router := NewPrefixRouter("/prefix")

	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, router.root, "Root should not be nil")
	assert.Equal(t, "/prefix", router.prefix, "Root should not be nil")

}
func TestRouter_RegisterSimplePath(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	r.Register(GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}

func TestRouter_RegisterNilHandler(t *testing.T) {
	r := NewRouter()

	assert.PanicsWithValue(t, "handler must not be nil", func() {
		r.Register(GET, "/path", nil)
	})
}
func TestRouter_NotFound(t *testing.T) {
	r := NewRouter()
	path := "/not-found"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	r.NotFound(method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_RegisterWithMiddleware(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware", "HIT")
			next.ServeHTTP(w, r)
		})
	}
	r.Use(middleware)
	r.Register(GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
	assert.Equal(t, "HIT", res.Header.Get("Middleware"))
}
func TestRouter_RegisterBeforeMiddleware(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Middleware", "HIT")
			next.ServeHTTP(w, r)
		})
	}
	r.Register(GET, path, method)

	assert.Panics(t, func() { r.Use(middleware) }, "Use should panic after declaring first route")

}
func TestRouter_RegisterCapturePathVariable(t *testing.T) {
	r := NewRouter()
	path := "/path/{test}"
	pathToTest := "/path/hello"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(context.Param(r, "test")))
	}
	r.Register(GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, pathToTest))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello", string(body))
}
func TestRouter_RegisterCaptureMultiplePathVariable(t *testing.T) {
	r := NewRouter()
	path := "/path/{test}/path/{test2}"
	pathToTest := "/path/hello/path/world"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s_%s", context.Param(r, "test"), context.Param(r, "test2"))))
	}
	r.Register(GET, path, method)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, pathToTest))

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
}
func TestRouter_RegisterMultiplePath(t *testing.T) {
	r := NewRouter()
	path1 := "/path/"
	path2 := "/path/test/complex"
	method1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world_1"))
	}
	method2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world_2"))
	}
	r.Register(GET, path1, method1)
	r.Register(GET, path2, method2)
	s := setup(r)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path1))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world_1", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, path2))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, "hello_world_2", string(body))
}
func TestRouter_Group(t *testing.T) {
	router := NewRouter()
	path1 := "/path1"
	group := "/group"
	path2 := "/path2"
	method := func(w http.ResponseWriter, r *http.Request) {
	}
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			w.Write([]byte(fmt.Sprintf("Hello World Middleware 1")))
		})
	}
	router.Use(fn)
	router.Register(GET, path1, method)
	router.Group(group, func(r *Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				w.Write([]byte(fmt.Sprintf("Hello World Middleware 2")))
			})
		})
		r.Register(GET, path2, method)
	})
	s := setup(router)
	defer s.Close()

	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path1))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Hello World Middleware 1", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, group+path2))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, "Hello World Middleware 2Hello World Middleware 1", string(body))
}

func TestRouter_NestedGroupsComposePrefixes(t *testing.T) {
	r := NewRouter()
	r.Group("/api", func(api *Router) {
		api.Group("/v1", func(version *Router) {
			version.Register(GET, "/users", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("users"))
			})
		})
	})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "users", response.Body.String())
}

func TestRouter_WithPreservesPrefix(t *testing.T) {
	r := NewPrefixRouter("/api")
	r.With().Register(GET, "/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "ok", response.Body.String())
}

func TestRouter_SupportsConcurrentRegistrationAndServing(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/stable", func(w http.ResponseWriter, r *http.Request) {})

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		for i := 0; i < 100; i++ {
			r.Register(GET, fmt.Sprintf("/dynamic/%d", i), func(w http.ResponseWriter, r *http.Request) {})
		}
	}()
	go func() {
		defer waitGroup.Done()
		for i := 0; i < 100; i++ {
			request := httptest.NewRequest(http.MethodGet, "/stable", nil)
			r.ServeHTTP(httptest.NewRecorder(), request)
		}
	}()

	waitGroup.Wait()
}
func TestRouter_With(t *testing.T) {
	r := NewRouter()
	path := "/path"
	method := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	}
	methodWith := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test_with"))
			next.ServeHTTP(w, r)
		})
	}
	r.Register(GET, path, method)
	s := setup(r)
	defer s.Close()

	r.With(methodWith).Register(GET, "/path_with", method)
	res, _ := http.Get(fmt.Sprintf("%s%s", s.URL, path))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "hello_world", string(body))
	res, _ = http.Get(fmt.Sprintf("%s%s", s.URL, "/path_with"))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, "test_withhello_world", string(body))
}

func TestRouter_WithInheritsParentMiddleware(t *testing.T) {
	r := NewRouter()
	tag := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(name))
				next.ServeHTTP(w, r)
			})
		}
	}
	r.Use(tag("parent"))
	r.With(tag("with")).Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handler"))
	})

	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, "parentwithhandler", response.Body.String())
}

func TestRouter_GroupRejectsMiddlewareAfterItsFirstRoute(t *testing.T) {
	r := NewRouter()
	group := r.Group("/group", func(subrouter *Router) {
		subrouter.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})
	})

	assert.Panics(t, func() {
		group.Use(func(next http.Handler) http.Handler { return next })
	})
}

func TestRouter_GroupAcceptsMiddlewareAfterSiblingRoute(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/ping", func(w http.ResponseWriter, r *http.Request) {})

	assert.NotPanics(t, func() {
		r.Group("/group", func(subrouter *Router) {
			subrouter.Use(func(next http.Handler) http.Handler { return next })
			subrouter.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})
		})
	})
}

func TestRouter_GroupNotFoundReachesServingRouter(t *testing.T) {
	r := NewRouter()
	r.Group("/group", func(subrouter *Router) {
		subrouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("custom"))
		})
		subrouter.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})
	})

	request := httptest.NewRequest(http.MethodGet, "/group/missing", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "custom", response.Body.String())
}

func TestRouter_NotFoundNilHandler(t *testing.T) {
	r := NewRouter()

	assert.PanicsWithValue(t, "handler must not be nil", func() { r.NotFound(nil) })
}

func TestRouter_TreatsParameterSyntaxInRequestPathAsLiteral(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/{id}", func(w http.ResponseWriter, r *http.Request) {
		routerCtx, _ := context.FromRequest(r)
		w.Write([]byte(routerCtx.Value("id")))
	})

	request := httptest.NewRequest(http.MethodGet, "/%7B*%7D", nil)
	response := httptest.NewRecorder()

	assert.NotPanics(t, func() { r.ServeHTTP(response, request) })
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{*}", response.Body.String())
}

func TestRouter_GroupResultIsUsableAfterTheCallback(t *testing.T) {
	r := NewRouter()
	group := r.Group("/group", func(subrouter *Router) {
		subrouter.Register(GET, "/first", func(w http.ResponseWriter, r *http.Request) {})
	})
	group.Register(GET, "/second", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("second"))
	})

	var handler http.Handler = group
	request := httptest.NewRequest(http.MethodGet, "/group/second", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "second", response.Body.String())
}

func TestRouter_RouteWithoutVariablesCarriesNoContext(t *testing.T) {
	r := NewRouter()
	var found bool
	var value string
	r.Register(GET, "/static", func(w http.ResponseWriter, r *http.Request) {
		routerCtx, ok := context.FromRequest(r)
		found = ok
		value = routerCtx.Value("anything")
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/static", nil))

	assert.False(t, found, "a route capturing nothing should not derive a request")
	assert.Equal(t, "", value, "reading through the absent context must stay safe")
}

func TestRouter_MethodNotAllowedHandlerIsConfigurable(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("custom"))
	})

	request := httptest.NewRequest(http.MethodPost, "/path", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusTeapot, response.Code)
	assert.Equal(t, "custom", response.Body.String())
	assert.Equal(t, "GET", response.Header().Get("Allow"), "Allow is set before the handler runs")
}

func TestRouter_GroupMethodNotAllowedReachesServingRouter(t *testing.T) {
	r := NewRouter()
	r.Group("/group", func(subrouter *Router) {
		subrouter.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})
		subrouter.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})
	})

	request := httptest.NewRequest(http.MethodPost, "/group/path", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusTeapot, response.Code)
}

func TestRouter_MethodNotAllowedNilHandler(t *testing.T) {
	r := NewRouter()

	assert.PanicsWithValue(t, "handler must not be nil", func() { r.MethodNotAllowed(nil) })
}

func TestRouter_RegisterPathWithQueryString(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello_world"))
	})

	request := httptest.NewRequest(http.MethodGet, "/path?query=value", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "hello_world", response.Body.String())
}

func TestRouter_RegisterRootPreservesExistingRoutes(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("path"))
	})
	r.Register(GET, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	})

	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "path", response.Body.String())
}

func TestRouter_ReturnsMethodNotAllowed(t *testing.T) {
	r := NewRouter()
	r.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})

	request := httptest.NewRequest(http.MethodPost, "/path", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	assert.Equal(t, "GET", response.Header().Get("Allow"))
}

func TestRouter_MethodNotAllowedListsEveryRegisteredMethod(t *testing.T) {
	r := NewRouter()
	r.Register(PATCH, "/path", func(w http.ResponseWriter, r *http.Request) {})
	r.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {})

	request := httptest.NewRequest(http.MethodDelete, "/path", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	assert.Equal(t, "GET, PATCH", response.Header().Get("Allow"))
}

func TestRouter_MethodNotAllowedSpansEveryMatchingBranch(t *testing.T) {
	r := NewRouter()
	r.Register(POST, "/a/b", func(w http.ResponseWriter, r *http.Request) {})
	r.Register(PATCH, "/a/{id}", func(w http.ResponseWriter, r *http.Request) {})

	request := httptest.NewRequest(http.MethodGet, "/a/b", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	assert.Equal(t, "PATCH, POST", response.Header().Get("Allow"))
}

func TestRouter_BacktracksWhenStaticMatchLacksMethod(t *testing.T) {
	r := NewRouter()
	r.Register(POST, "/a/b", func(w http.ResponseWriter, r *http.Request) {})
	r.Register(GET, "/a/{id}", func(w http.ResponseWriter, r *http.Request) {
		routerCtx, _ := context.FromRequest(r)
		w.Write([]byte(routerCtx.Value("id")))
	})

	request := httptest.NewRequest(http.MethodGet, "/a/b", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "b", response.Body.String())
}

func TestRouter_ReturnsMethodNotAllowedWhenNoBranchMatchesMethod(t *testing.T) {
	r := NewRouter()
	r.Register(POST, "/a/b", func(w http.ResponseWriter, r *http.Request) {})
	r.Register(PATCH, "/a/{id}", func(w http.ResponseWriter, r *http.Request) {})

	request := httptest.NewRequest(http.MethodGet, "/a/b", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
}

func TestRouter_MiddlewareOrder(t *testing.T) {
	r := NewRouter()
	order := make([]string, 0, 3)
	appendMiddleware := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	r.Use(appendMiddleware("first"))
	r.Use(appendMiddleware("second"))
	r.Register(GET, "/path", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	r.ServeHTTP(httptest.NewRecorder(), request)

	assert.Equal(t, []string{"first", "second", "handler"}, order)
}

func setup(r http.Handler) *httptest.Server {

	return httptest.NewServer(r)
}
