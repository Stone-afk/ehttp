package accesslog

import (
	"testing"
	"time"
	v6 "web/v6"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	logbd := NewBuilder()
	s := v6.NewHTTPServer()
	s.Get("/", func(ctx *v6.Context) {
		ctx.Response.Write([]byte("hello, world"))
	})

	s.Get("/user", func(ctx *v6.Context) {
		time.Sleep(time.Second)
		ctx.RespData = []byte("hello, user")
	})
	s.Use(logbd.Build())
	s.Start(":8081")
}
