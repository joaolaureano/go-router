package routing

import (
	"slices"
	"strings"
)

// node is one edge of a compressed (patricia) trie. path is the literal text
// this edge adds relative to its parent -- not necessarily a whole path
// segment, and not necessarily bounded by a single one either: two routes that
// share a long literal run, even across a "/", share one chain of edges for it,
// split apart only at the byte where they first diverge.
//
// Only a parameter placeholder forces a node boundary, because only it needs
// one: a placeholder always spans exactly one whole segment, so it can only
// ever anchor to a node that represents having fully consumed some number of
// whole segments. Nothing else about where an edge starts or ends carries
// meaning, which is what lets registration split and merge edges freely
// without disturbing that invariant.
type node[E any] struct {
	path      string
	children  []*node[E]
	parameter *node[E]
	endpoints []endpoint[E]
}

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

// linearScanLimit is where scanning stops beating bisection. Below it the
// branch-free walk over a handful of pointers wins; above it the comparisons
// saved are worth the mispredictions.
const linearScanLimit = 8

// childByFirstByte returns the child whose edge could match b, and the index
// it occupies (or where it would belong, if nil is returned).
//
// A byte, not the whole edge, is enough: children are kept disjoint on their
// first byte by construction (ensureStaticChild and mergeStaticChild never
// leave two children agreeing there), so at most one can ever answer, and this
// is the dispatch every walk and every insertion is built on.
//
// children is kept sorted so that a node with many of them costs log(n)
// comparisons rather than n. A router fanning out to fifty resources at one
// level is ordinary, and scanning them all was the largest cost in the walk
// by some way.
func (n *node[E]) childByFirstByte(b byte) (int, *node[E]) {
	children := n.children
	if len(children) < linearScanLimit {
		for i, child := range children {
			if child.path[0] == b {
				return i, child
			}
			if child.path[0] > b {
				return i, nil
			}
		}
		return len(children), nil
	}

	low, high := 0, len(children)
	for low < high {
		middle := int(uint(low+high) >> 1)
		if children[middle].path[0] < b {
			low = middle + 1
		} else {
			high = middle
		}
	}
	if low < len(children) && children[low].path[0] == b {
		return low, children[low]
	}
	return low, nil
}

// insertChildAt splices child in at the position childByFirstByte named, kept
// sorted so that childByFirstByte can bisect. Registration pays the shift; a
// lookup would pay for the scan on every request.
func (n *node[E]) insertChildAt(at int, child *node[E]) {
	n.children = append(n.children, nil)
	copy(n.children[at+1:], n.children[at:])
	n.children[at] = child
}

// commonPrefixLen returns how much of a and b agree from the start.
func commonPrefixLen(a, b string) int {
	limit := min(len(a), len(b))
	i := 0
	for i < limit && a[i] == b[i] {
		i++
	}
	return i
}

// splitEdgeAt breaks the edge at children[at] into two: an intermediate node
// holding the first commonLen bytes, with the original -- now shortened to
// what follows -- as its sole child. It returns the intermediate node, which
// is where whatever needs the shared prefix now belongs.
func splitEdgeAt[E any](target *node[E], at int, commonLen int) *node[E] {
	existing := target.children[at]
	split := newNode[E](existing.path[:commonLen])
	existing.path = existing.path[commonLen:]
	split.children = []*node[E]{existing}
	target.children[at] = split
	return split
}

// ensureStaticChild finds or creates the node representing key's literal text
// under target, splitting an existing edge wherever key and it first diverge.
// key never crosses into a placeholder: it is exactly the literal run a
// pattern declares between two placeholders (or between one and either end of
// the pattern), which is what lets a placeholder always anchor to the node
// this returns.
func ensureStaticChild[E any](target *node[E], key string) *node[E] {
	if key == "" {
		return target
	}

	at, existing := target.childByFirstByte(key[0])
	if existing == nil {
		leaf := newNode[E](key)
		target.insertChildAt(at, leaf)
		return leaf
	}

	common := commonPrefixLen(existing.path, key)
	if common < len(existing.path) {
		existing = splitEdgeAt(target, at, common)
	}
	return ensureStaticChild(existing, key[common:])
}

// setEndpoint records what this node resolves to under a method. It only ever
// appends: registering a method twice on one path is rejected before the walk
// gets here, and a merge asks whether the method is taken before grafting.
func (n *node[E]) setEndpoint(httpMethod Method, handler E, variableNames []string) {
	n.endpoints = append(n.endpoints, endpoint[E]{
		method:        httpMethod,
		handler:       handler,
		variableNames: variableNames,
	})
}

func (n *node[E]) hasMethod(httpMethod Method) bool {
	return n.find(httpMethod) != nil
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

	// allowed accumulates the methods a path answers under when the walk reaches
	// it with the wrong one. A slice, not a set: a path is registered under one
	// or two methods and never more than the seven that exist, and at that size
	// scanning beats hashing -- and the slice is the very one handed back, where
	// a set had to be copied out into one to be sorted.
	allowed []Method
}

// typicalCaptureCount sizes the slice on first capture. Growing from nothing
// costs an allocation per variable; starting at two spares that for the routes
// people actually write, and costs a route with a single variable nothing but a
// slot it does not use. Anything wider only pays off on patterns rare enough
// not to size for.
const typicalCaptureCount = 2

// capture records a value the walk matched against a parameter.
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
	slices.Sort(search.allowed)
	return search.allowed
}

// typicalMethodCount sizes the allowed slice on first use. Growing from nothing
// costs an allocation per verb recorded -- three of them on an ordinary CRUD
// path, where sizing once costs one. Six, not four, because the caller most
// likely to read this list is rendering an Allow header, and it appends the two
// verbs answered on the route's behalf before doing so.
const typicalMethodCount = 6

// recordAllowed notes the methods this node answers under. More than one node
// can reach here in a single walk -- a parameter branch and a static one both
// ending the same path -- so a method already recorded is skipped.
func (n *node[E]) recordAllowed(search *lookup) {
	if search.allowed == nil && len(n.endpoints) > 0 {
		search.allowed = make([]Method, 0, typicalMethodCount)
	}
	for i := range n.endpoints {
		if !slices.Contains(search.allowed, n.endpoints[i].method) {
			search.allowed = append(search.allowed, n.endpoints[i].method)
		}
	}
}

// splitSegment peels the first "/"-delimited segment off path, the way a raw
// request path is walked one segment at a time wherever decoding is involved.
func splitSegment(path string) (segment, rest string) {
	if separator := strings.IndexByte(path, '/'); separator >= 0 {
		return path[:separator], path[separator+1:]
	}
	return path, ""
}

// match walks what is left of the raw path, preferring the static edge and
// only then the parameter branch, and returns the first node answering the
// method being searched for.
//
// It reads the path with a cursor rather than a list of segments, so an
// ordinary lookup allocates nothing at all. An empty path means the walk has
// arrived: the caller trims the separators off the ends, so only the root
// reaches this with nothing left on the first call.
func (n *node[E]) match(path string, search *lookup) *node[E] {
	if path == "" {
		if n.hasMethod(search.httpMethod) {
			return n
		}
		n.recordAllowed(search)
		return nil
	}

	if _, child := n.childByFirstByte(path[0]); child != nil &&
		len(path) >= len(child.path) && path[:len(child.path)] == child.path {
		if matched := child.match(path[len(child.path):], search); matched != nil {
			return matched
		}
	}

	if n.parameter != nil {
		segment, rest := splitSegment(path)
		search.capture(segment)
		if matched := n.parameter.match(rest, search); matched != nil {
			return matched
		}
		search.uncapture()
	}

	return nil
}

// decodedCursor presents a raw path as decoded text, one segment at a time, so
// that matching against the tree's (always literal) edges can compare decoded
// bytes throughout -- without ever decoding more of the path than a match
// actually consumes.
//
// segment and raw are kept apart deliberately, rather than joined into one
// buffer: a "/" a segment's own decoding produces -- a %2F, the very thing
// that has to stay inside one variable -- would otherwise be indistinguishable
// from the separator between two segments, and a search for one would find
// the other instead.
//
// It is a plain value: passing it by copy to each recursive call gives every
// branch its own snapshot for free, exactly as passing path by value already
// does for the undecoded walk. A failed branch simply discards its copy;
// nothing needs to be saved and restored by hand.
type decodedCursor struct {
	segment string // decoded text of the raw segment currently being matched
	raw     string // whatever of the raw path lies beyond it, undecoded
	decode  func(string) string
}

// newDecodedCursor decodes the first segment of raw, so that a cursor is
// always ready to compare or capture without a special case for "nothing
// pulled yet".
func newDecodedCursor(raw string, decode func(string) string) decodedCursor {
	return decodedCursor{raw: raw, decode: decode}.pullSegment()
}

// pullSegment decodes the next raw segment into segment, discarding whatever
// was left of the one before it -- always called once that one is either
// fully matched or fully captured, never partway through.
func (c decodedCursor) pullSegment() decodedCursor {
	segment, rest := splitSegment(c.raw)
	return decodedCursor{segment: c.decode(segment), raw: rest, decode: c.decode}
}

// advance moves past the current segment entirely: to the next one if the raw
// path has one, or to the exhausted state match reads as "nothing left".
func (c decodedCursor) advance() decodedCursor {
	if c.raw == "" {
		return decodedCursor{decode: c.decode}
	}
	return c.pullSegment()
}

// exhausted reports whether every byte of the request has been matched or
// captured, with nothing decoded and nothing raw left to decode.
func (c decodedCursor) exhausted() bool {
	return c.segment == "" && c.raw == ""
}

// matchDecoded is match for a lookup that needs percent-decoding. Decoding
// happens lazily, one raw segment at a time, only as far as a candidate edge
// or a parameter capture actually requires -- the same walk as match, just
// reading through decodedCursor instead of the raw path directly.
func (n *node[E]) matchDecoded(cur decodedCursor, search *lookup) *node[E] {
	if cur.exhausted() {
		if n.hasMethod(search.httpMethod) {
			return n
		}
		n.recordAllowed(search)
		return nil
	}

	if cur.segment != "" {
		if child, next, ok := n.matchStaticDecoded(cur); ok {
			if matched := child.matchDecoded(next, search); matched != nil {
				return matched
			}
		}
	}

	if n.parameter != nil {
		search.capture(cur.segment)
		if matched := n.parameter.matchDecoded(cur.advance(), search); matched != nil {
			return matched
		}
		search.uncapture()
	}

	return nil
}

// matchStaticDecoded finds a child whose edge matches what cur has decoded so
// far, pulling and decoding one further raw segment at a time as the edge's
// own embedded "/" boundaries are reached. Segments are compared whole against
// the edge's own text rather than joined into a shared buffer first, which is
// what keeps a "/" a decoding produces from ever being mistaken for one of
// those boundaries.
func (n *node[E]) matchStaticDecoded(cur decodedCursor) (*node[E], decodedCursor, bool) {
	_, child := n.childByFirstByte(cur.segment[0])
	if child == nil {
		return nil, cur, false
	}

	edge := child.path
	for {
		switch {
		case len(cur.segment) < len(edge):
			// The segment is a strict prefix of what is left of the edge: the
			// edge can only continue with a literal "/", crossing into the
			// next raw segment, since edges are only ever merged at whole
			// segment boundaries.
			if edge[len(cur.segment)] != '/' || cur.segment != edge[:len(cur.segment)] {
				return nil, cur, false
			}
			if cur.raw == "" {
				return nil, cur, false
			}
			edge = edge[len(cur.segment)+1:]
			cur = cur.pullSegment()
		case len(cur.segment) == len(edge):
			if cur.segment != edge {
				return nil, cur, false
			}
			return child, cur.advance(), true
		default: // len(cur.segment) > len(edge)
			// The edge ends inside this segment's own decoded text: whatever
			// is left of it stays the current segment, for child's own
			// children to keep matching -- no more of the raw path is
			// consulted until this segment is spent.
			if cur.segment[:len(edge)] != edge {
				return nil, cur, false
			}
			cur.segment = cur.segment[len(edge):]
			return child, cur, true
		}
	}
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
		mergeStaticChild(target, sourceChild, prefixVariables)
	}

	mergeBranch(&target.parameter, source.parameter, prefixVariables)
}

// mergeStaticChild grafts source -- an edge and its subtree -- into target's
// children. The two sides may have split what is logically the same literal
// text at different points, so this splits and descends exactly as
// ensureStaticChild does, merging the two subtrees only once both edges agree
// all the way to one end of the shorter one.
func mergeStaticChild[E any](target *node[E], source *node[E], prefixVariables []string) {
	at, existing := target.childByFirstByte(source.path[0])
	if existing == nil {
		target.insertChildAt(at, cloneNode(source, prefixVariables))
		return
	}

	common := commonPrefixLen(existing.path, source.path)
	if common < len(existing.path) {
		existing = splitEdgeAt(target, at, common)
	}

	remainder := source.path[common:]
	if remainder == "" {
		mergeNodes(existing, source, prefixVariables)
		return
	}
	mergeStaticChild(existing, withPath(source, remainder), prefixVariables)
}

// withPath returns a shallow, read-only view of n under a different path. It
// exists so that mergeStaticChild can recurse with "the part of source not yet
// matched" without mutating source, which belongs to the other tree.
func withPath[E any](n *node[E], path string) *node[E] {
	return &node[E]{
		path:      path,
		children:  n.children,
		parameter: n.parameter,
		endpoints: n.endpoints,
	}
}

// mergeBranch merges the parameter branch, cloning when the target has nothing
// there yet.
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
	// source.children is already sorted, so copying in order preserves it.
	for _, child := range source.children {
		clone.children = append(clone.children, cloneNode(child, prefixVariables))
	}
	if source.parameter != nil {
		clone.parameter = cloneNode(source.parameter, prefixVariables)
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
