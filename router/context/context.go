package context

import "slices"

const ContextKey = "RouterContext"

type Context struct {
	paramNames []string
	paramValue []string
}

func (rp *Context) Value(key string) string {
	idx := slices.Index(rp.paramNames, key)
	if idx > -1 {
		return rp.paramValue[idx]
	}
	return ""
}

func (rp *Context) Set(key string, value string) {
	rp.paramNames = append(rp.paramNames, key)
	rp.paramValue = append(rp.paramValue, value)
}

func NewContext() *Context {
	return &Context{
		paramNames: make([]string, 0),
		paramValue: make([]string, 0),
	}
}
