package v4

import "net/http"

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	PathParams map[string]string
}
