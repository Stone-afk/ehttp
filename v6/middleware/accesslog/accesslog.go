package accesslog

import (
	"encoding/json"
	"log"
	v6 "web/v6"
)

func (b *MiddlewareBuilder) Build() v6.Middleware {
	return func(next v6.HandleFunc) v6.HandleFunc {
		return func(ctx *v6.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Request.Host,
					Path:       ctx.Request.URL.Path,
					HTTPMethod: ctx.Request.Method,
					Route:      ctx.MatchedRoute,
				}
				val, _ := json.Marshal(l)
				b.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

func (b *MiddlewareBuilder) LogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}
