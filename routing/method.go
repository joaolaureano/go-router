package routing

// Method is an HTTP method as the routing table understands it. It stays a
// string underneath so that a verb the constants below do not name still
// routes; the constants exist to spell the common ones without typos.
//
// Spelled out rather than taken from net/http so that this package keeps to
// its own vocabulary. The values are fixed by RFC 9110 and cannot drift.
type Method string

const (
	GET     Method = "GET"
	HEAD    Method = "HEAD"
	POST    Method = "POST"
	PUT     Method = "PUT"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
)
