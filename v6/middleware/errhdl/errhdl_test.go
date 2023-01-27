package errhdl

import (
	"bytes"
	"html/template"
	"testing"
	v6 "web/v6"
)

func TestNewMiddlewareBuilder(t *testing.T) {
	s := v6.NewHTTPServer()
	s.Get("/user", func(ctx *v6.Context) {
		ctx.RespData = []byte("hello, world")
	})
	page := `
<html>
	<h1>404 NOT FOUND 我的自定义错误页面</h1>
</html>
`
	tpl, err := template.New("404").Parse(page)
	if err != nil {
		t.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.Use(NewMiddlewareBuilder().
		RegisterError(404, buffer.Bytes()).Build())

	s.Start(":8081")
}
