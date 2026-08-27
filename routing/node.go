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
	endpoints map[Method]endpoint[E]
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
type endpoint[E any] struct {
	handler       E
	variableNames []string
}

func newNode[E any](path string) *node[E] {
	return &node[E]{
		path:      path,
		children:  make([]*node[E], 0),
		endpoints: make(map[Method]endpoint[E]),
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
	n.endpoints[httpMethod] = endpoint[E]{
		handler:       handler,
		variableNames: variableNames,
	}
}

func (n *node[E]) hasMethod(httpMethod Method) bool {
	_, exists := n.endpoints[httpMethod]
	return exists
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

	allowed map[Method]struct{}
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
	for httpMethod := range n.endpoints {
		search.allowed[httpMethod] = struct{}{}
	}
}

// match walks the remaining segments, preferring the static child, then the
// parameter branch, and only then the catch-all, and returns the first node
// answering the method being searched for.
//
// The catch-all is also tried once the path runs out, so that "/files/*" covers
// "/files" itself with an empty remainder.
func (n *node[E]) match(paths []string, index int, search *lookup) *node[E] {
	if index == len(paths) {
		if n.hasMethod(search.httpMethod) {
			return n
		}
		n.recordAllowed(search)
		return n.matchWildcard(paths, index, search)
	}

	if child := n.staticChild(paths[index]); child != nil {
		if matched := child.match(paths, index+1, search); matched != nil {
			return matched
		}
	}

	if n.parameter != nil {
		search.params = append(search.params, Param{Value: paths[index]})
		if matched := n.parameter.match(paths, index+1, search); matched != nil {
			return matched
		}
		search.params = search.params[:len(search.params)-1]
	}

	return n.matchWildcard(paths, index, search)
}

// matchWildcard consumes whatever is left of the path in one value. A catch-all
// node is always terminal, so there is nothing further to walk.
func (n *node[E]) matchWildcard(paths []string, index int, search *lookup) *node[E] {
	if n.wildcard == nil {
		return nil
	}

	search.params = append(search.params, Param{Value: strings.Join(paths[index:], "/")})
	if n.wildcard.hasMethod(search.httpMethod) {
		return n.wildcard
	}
	n.wildcard.recordAllowed(search)
	search.params = search.params[:len(search.params)-1]
	return nil
}

// mergeNodes copies source into target. prefixVariables names the variables
// that target already sits below, which every grafted endpoint has to inherit.
func mergeNodes[E any](target, source *node[E], prefixVariables []string) {
	for httpMethod, sourceEndpoint := range source.endpoints {
		if !target.hasMethod(httpMethod) {
			target.endpoints[httpMethod] = cloneEndpoint(sourceEndpoint, prefixVariables)
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
	for httpMethod, sourceEndpoint := range source.endpoints {
		clone.endpoints[httpMethod] = cloneEndpoint(sourceEndpoint, prefixVariables)
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
		handler:       source.handler,
		variableNames: variableNames,
	}
}
