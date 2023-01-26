package v4

import "testing"

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, world"))
	})

	s.Get("/user", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, user"))
	})

	s.Start("127.0.0.1:8090")

}
