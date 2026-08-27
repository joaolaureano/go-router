package routing

import (
	"sort"
	"strings"
)

// node is one path segment of the routing. It is deliberately unexported: the
// invariant tying an endpoint's variable names to the parameter nodes above it
// only holds while Tree is the sole writer.
type node[E any] struct {
	path      string
	children  []*node[E]
	parameter *node[E]
	wildcard  *node[E]
	endpoints []endpoint[E]
}

// segmentKind tells the three shapes a pattern segment can take apart. Only the
// registrar classifies a segment: to a request, every segment is literal text.
type segmentKind int

const (
	staticSegment segmentKind = iota
	parameterSegment
	wildcardSegment
)

// endpoint is what a single method registered on a node resolves to. The
// variable names live here rather than on the node because two methods on the
// same path may have been written with different parameter names.
//
// Nodes hold these in a slice, not a map keyed by method. A path answers one or
// two methods in practice, and at that size scanning a contiguous slice beats
// hashing a string. It also means an intermediate node -- most of the tree --
// allocates nothing at all, where a map header cost every node 48 bytes whether
// or not a route ever ended there.
type endpoint[E any] struct {
	method        Method
	handler       E
	variableNames []string
}

// find returns the endpoint registered for this method, or nil.
func (n *node[E]) find(httpMethod Method) *endpoint[E] {
	for i := range n.endpoints {
		if n.endpoints[i].method == httpMethod {
			return &n.endpoints[i]
		}
	}
	return nil
}

func newNode[E any](path string) *node[E] {
	return &node[E]{
		path:     path,
		children: make([]*node[E], 0),
	}
}

// staticChild finds the child holding this exact segment. It never returns the
// parameter or catch-all branch: a request segment that happens to read "{*}"
// or "*" is literal text, not a request for those slots.
func (n *node[E]) staticChild(path string) *node[E] {
	for _, child := range n.children {
		if path == child.path {
			return child
		}
	}
	return nil
}

// childFor returns the slot a registration should descend into. Only the
// registrar knows how a segment was written, so only it may ask for the
// parameter or wildcard branch.
func (n *node[E]) childFor(path string, kind segmentKind) *node[E] {
	switch kind {
	case parameterSegment:
		return n.parameter
	case wildcardSegment:
		return n.wildcard
	default:
		return n.staticChild(path)
	}
}

func (n *node[E]) addChild(child *node[E], kind segmentKind) {
	if child == nil {
		panic("node child must not be nil")
	}
	switch kind {
	case parameterSegment:
		n.parameter = child
	case wildcardSegment:
		n.wildcard = child
	default:
		n.children = append(n.children, child)
	}
}

func (n *node[E]) setEndpoint(httpMethod Method, handler E, variableNames []string) {
	if existing := n.find(httpMethod); existing != nil {
		existing.handler = handler
		existing.variableNames = variableNames
		return
	}
	n.endpoints = append(n.endpoints, endpoint[E]{
		method:        httpMethod,
		handler:       handler,
		variableNames: variableNames,
	})
}

func (n *node[E]) hasMethod(httpMethod Method) bool {
	return n.find(httpMethod) != nil
}

func (n *node[E]) hasAnyMethod() bool {
	return len(n.endpoints) > 0
}

// lookup carries the state of one walk. Parameter values accumulate as the
// walk descends and unwind when a branch fails, so only the values on the
// surviving route remain. Nodes that end the path under some other verb are
// recorded on the way, which lets a single walk answer both "which handler"
// and "which methods would have worked".
//
// It is not parameterised on the endpoint type: nothing it holds depends on
// what a route resolves to.
type lookup struct {
	httpMethod Method

	// params accumulate as the walk descends and unwind when a branch fails, so
	// only those on the surviving route remain. Values land here as they are
	// captured and the names are filled in at the end, once the matched
	// endpoint says what it called them -- one slice built once, rather than a
	// slice of values and a second one to pair it with.
	params []Param

	// decode transforms each segment as the walk reaches it, or is nil when the
	// path needs nothing done to it -- which is the usual case, and the reason
	// the walk can run straight off the request path without splitting it.
	decode func(string) string

	allowed map[Method]struct{}
}

// segmentEnd reports where the leading segment ends: at the next separator, or
// at the end of the path when there is none left.
//
// It returns an index rather than the two substrings so that it stays inside
// the inliner's budget -- constructing the strings here costs enough to push it
// over, and this runs once per segment of every request.
func segmentEnd(path string) int {
	if separator := strings.IndexByte(path, '/'); separator >= 0 {
		return separator
	}
	return len(path)
}

// typicalCaptureCount sizes the slice on first capture. Growing from nothing
// costs an allocation per variable; starting at two spares that for the routes
// people actually write, and costs a route with a single variable nothing but a
// slot it does not use. Anything wider only pays off on patterns rare enough
// not to size for.
const typicalCaptureCount = 2

// capture records a value the walk matched against a parameter or catch-all.
// The name stays blank until the matched endpoint says what it is called.
func (search *lookup) capture(value string) {
	if search.params == nil {
		search.params = make([]Param, 0, typicalCaptureCount)
	}
	search.params = append(search.params, Param{Value: value})
}

// uncapture drops the value a branch recorded before it turned out not to
// match, so only the surviving route's variables remain.
func (search *lookup) uncapture() {
	search.params = search.params[:len(search.params)-1]
}

// allowedMethods returns the recorded methods sorted, so that a caller
// rendering them -- into an Allow header, say -- gets a stable order.
func (search *lookup) allowedMethods() []Method {
	if len(search.allowed) == 0 {
		return nil
	}
	methods := make([]Method, 0, len(search.allowed))
	for httpMethod := range search.allowed {
		methods = append(methods, httpMethod)
	}
	sort.Slice(methods, func(i, j int) bool { return methods[i] < methods[j] })
	return methods
}

func (n *node[E]) recordAllowed(search *lookup) {
	if !n.hasAnyMethod() {
		return
	}
	if search.allowed == nil {
		search.allowed = make(map[Method]struct{}, len(n.endpoints))
	}
	for i := range n.endpoints {
		search.allowed[n.endpoints[i].method] = struct{}{}
	}
}

// match walks what is left of the path, preferring the static child, then the
// parameter branch, and only then the catch-all, and returns the first node
// answering the method being searched for.
//
// It reads the path with a cursor rather than a list of segments, so an
// ordinary lookup allocates nothing at all. An empty path means the walk has
// arrived: the caller trims the separators off the ends, so only the root
// reaches this with nothing left on the first call.
//
// The catch-all is tried once the path runs out too, so that "/files/*" covers
// "/files" itself with an empty remainder.
func (n *node[E]) match(path string, search *lookup) *node[E] {
	if path == "" {
		if n.hasMethod(search.httpMethod) {
			return n
		}
		n.recordAllowed(search)
		return n.matchWildcard(path, search)
	}

	// IndexByte is called here rather than behind a helper: it costs enough on
	// its own that any wrapper around it exceeds the inliner's budget, and this
	// runs once per segment of every request.
	segment, rest := path, ""
	if separator := strings.IndexByte(path, '/'); separator >= 0 {
		segment, rest = path[:separator], path[separator+1:]
	}
	if search.decode != nil {
		segment = search.decode(segment)
	}

	if child := n.staticChild(segment); child != nil {
		if matched := child.match(rest, search); matched != nil {
			return matched
		}
	}

	if n.parameter != nil {
		search.capture(segment)
		if matched := n.parameter.match(rest, search); matched != nil {
			return matched
		}
		search.uncapture()
	}

	return n.matchWildcard(path, search)
}

// matchWildcard consumes whatever is left of the path in one value. A catch-all
// node is always terminal, so there is nothing further to walk -- and since the
// remainder is a slice of the path already, taking it costs nothing.
func (n *node[E]) matchWildcard(path string, search *lookup) *node[E] {
	if n.wildcard == nil {
		return nil
	}

	remainder := path
	if search.decode != nil {
		remainder = search.decode(remainder)
	}

	search.capture(remainder)
	if n.wildcard.hasMethod(search.httpMethod) {
		return n.wildcard
	}
	n.wildcard.recordAllowed(search)
	search.uncapture()
	return nil
}

// mergeNodes copies source into target. prefixVariables names the variables
// that target already sits below, which every grafted endpoint has to inherit.
func mergeNodes[E any](target, source *node[E], prefixVariables []string) {
	for i := range source.endpoints {
		if !target.hasMethod(source.endpoints[i].method) {
			target.endpoints = append(target.endpoints, cloneEndpoint(source.endpoints[i], prefixVariables))
		}
	}

	for _, sourceChild := range source.children {
		targetChild := target.staticChild(sourceChild.path)
		if targetChild == nil {
			target.children = append(target.children, cloneNode(sourceChild, prefixVariables))
			continue
		}
		mergeNodes(targetChild, sourceChild, prefixVariables)
	}

	mergeBranch(&target.parameter, source.parameter, prefixVariables)
	mergeBranch(&target.wildcard, source.wildcard, prefixVariables)
}

// mergeBranch merges one of the single-slot branches, cloning when the target
// has nothing there yet.
func mergeBranch[E any](target **node[E], source *node[E], prefixVariables []string) {
	if source == nil {
		return
	}
	if *target == nil {
		*target = cloneNode(source, prefixVariables)
		return
	}
	mergeNodes(*target, source, prefixVariables)
}

func cloneNode[E any](source *node[E], prefixVariables []string) *node[E] {
	clone := newNode[E](source.path)
	for i := range source.endpoints {
		clone.endpoints = append(clone.endpoints, cloneEndpoint(source.endpoints[i], prefixVariables))
	}
	for _, child := range source.children {
		clone.children = append(clone.children, cloneNode(child, prefixVariables))
	}
	if source.parameter != nil {
		clone.parameter = cloneNode(source.parameter, prefixVariables)
	}
	if source.wildcard != nil {
		clone.wildcard = cloneNode(source.wildcard, prefixVariables)
	}
	return clone
}

func cloneEndpoint[E any](source endpoint[E], prefixVariables []string) endpoint[E] {
	variableNames := make([]string, 0, len(prefixVariables)+len(source.variableNames))
	variableNames = append(variableNames, prefixVariables...)
	variableNames = append(variableNames, source.variableNames...)
	return endpoint[E]{
		method:        source.method,
		handler:       source.handler,
		variableNames: variableNames,
	}
}
