package v5

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, world"))
	})

	s.Get("/user", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, user"))
	})

	s.Get("/user/:id", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, user param"))
	})

	s.Get("/a/b/*", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, a,b start"))
	})

	s.Get("/order/*", func(ctx *Context) {
		ctx.Response.Write([]byte("hello, order start"))
	})

	// 正则匹配
	s.Get("/sku/:id(^[0-9]+$)", func(ctx *Context) {
		ctx.Response.Write([]byte("hello,regx route"))
	})

	s.Start("127.0.0.1:8090")

}
