package _const

import "net/http"

type HTTPMethods string

const (
	GET    = http.MethodGet
	POST   = http.MethodPost
	PATCH  = http.MethodPatch
	DELETE = http.MethodDelete
)
