package v2

import "net/http"

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}
