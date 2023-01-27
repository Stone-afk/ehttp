package v6

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

func (ctx *Context) RespJSONOK(val any) error {
	return ctx.RespJSON(http.StatusOK, val)
}

func (ctx *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}

	//ctx.Response.WriteHeader(code)
	//_, err = ctx.Response.Write(bs)
	ctx.RespStatusCode = code
	ctx.RespData = bs
	return err
}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.Response, cookie)
}

func (ctx *Context) PathValue(key string) StringValue {
	val, ok := ctx.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: val}
}

func (ctx *Context) QueryValue(key string) StringValue {
	if ctx.cacheQueryValues == nil {
		ctx.cacheQueryValues = ctx.Request.URL.Query()
	}
	vals, ok := ctx.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: vals[0]}

}

func (ctx *Context) FormValue(key string) StringValue {
	if err := ctx.Request.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: ctx.Request.FormValue(key)}
}

func (ctx *Context) BindJSON(val any) error {
	if ctx.Request.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

type StringValue struct {
	val string
	err error
}

type Context struct {
	Request *http.Request
	// Response 原生的 ResponseWriter。当你直接使用 Response 的时候，
	// 那么相当于你绕开了 RespStatusCode 和 RespData。
	// 响应数据直接被发送到前端，其它中间件将无法修改响应
	// 其实我们也可以考虑将这个做成私有的
	Response http.ResponseWriter

	// 缓存的响应部分
	// 这部分数据会在最后刷新
	RespStatusCode int
	RespData       []byte

	// 路径参数
	PathParams map[string]string

	// 命中的路由
	MatchedRoute string

	// 万一将来有需求，可以考虑支持这个，但是需要复杂一点的机制
	// Body []byte 用户返回的响应
	// Err error 用户执行的 Error

	// 缓存的数据
	cacheQueryValues url.Values
}
