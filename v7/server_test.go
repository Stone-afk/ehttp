package v7

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
	"web/v7/middleware/accesslog"
	"web/v7/middleware/errhdl"
	"web/v7/middleware/opentelemetry"
	"web/v7/middleware/recovery"
)

func login(ctx *Context) {
	tpl := template.New("login")
	tpl, err := tpl.Parse(`
<html>
	<body>
		<form>
			// 在这里继续写页面
		<form>
	</body>
</html>
`)
	if err != nil {
		fmt.Println(err)
	}
	page := &bytes.Buffer{}
	err = tpl.Execute(page, nil)
	if err != nil {
		fmt.Println(err)
	}
	ctx.RespStatusCode = 200
	ctx.RespData = page.Bytes()
}

func TestServer(t *testing.T) {

	s := NewHTTPServer(ServerWithTemplateEngine(&GoTemplateEngine{}))

	err := s.TplEngine.LoadFromGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}

	logBd := accesslog.NewBuilder()
	bufByte, err := s.TplEngine.ExcuteTpl()
	if err != nil {
		t.Fatal(err)
	}
	erHdl := errhdl.NewBuilder().RegisterError(404, bufByte)
	tracBd := &opentelemetry.MiddlewareBuilder{}
	pacnicBd := &recovery.MiddlewareBuilder{}

	s.Use(erHdl.Build(), tracBd.Build(), logBd.Build(), pacnicBd.Build())

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

	// 模板渲染
	s.Get("login", login)

	s.Start("127.0.0.1:8090")

}

//func TestServerWithRenderEngine(t *testing.T) {
//
//	s := NewHTTPServer(ServerWithTemplateEngine(&GoTemplateEngine{}))
//
//	err := s.tplEngine.LoadFromGlob("testdata/tpls/*.gohtml")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	s.Get("/login", func(ctx *Context) {
//		er := ctx.Render("login.gohtml", nil)
//		if er != nil {
//			t.Fatal(er)
//		}
//	})
//
//	err = s.Start(":8081")
//
//	if err != nil {
//		t.Fatal(err)
//	}
//}
