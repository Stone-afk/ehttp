package v5

import "net/http"

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	PathParams map[string]string
}
