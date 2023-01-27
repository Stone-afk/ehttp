package repeat_body

import (
	"io/ioutil"
	web "web/v8"
)

func Middleware() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			ctx.Request.Body = ioutil.NopCloser(ctx.Request.Body)
			next(ctx)
		}
	}
}
