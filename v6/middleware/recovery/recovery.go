package recovery

import v6 "web/v6"

func (m *MiddlewareBuilder) Build() v6.Middleware {
	return func(next v6.HandleFunc) v6.HandleFunc {
		return func(ctx *v6.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = []byte(m.ErrMsg)
					// 万一 LogFunc 也panic，那我们也无能为力了
					m.LogFunc(ctx)
				}
			}()
			// 这里就是before route, before execute
			next(ctx)
			// 这里就是after route, after execute
		}
	}
}

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *v6.Context)
}
