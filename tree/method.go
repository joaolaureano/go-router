package tree

import "net/http"

// Method is an HTTP method as the routing table understands it. It stays a
// string underneath so that a verb the constants below do not name still
// routes; the constants exist to spell the common ones without typos.
type Method string

const (
	GET     Method = http.MethodGet
	HEAD    Method = http.MethodHead
	POST    Method = http.MethodPost
	PUT     Method = http.MethodPut
	PATCH   Method = http.MethodPatch
	DELETE  Method = http.MethodDelete
	OPTIONS Method = http.MethodOptions
)
