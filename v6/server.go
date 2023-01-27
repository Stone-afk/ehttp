package v6

import (
	"log"
	"net"
	"net/http"
)

const (
	zero     = 0
	notFound = 404
)

func (s *HTTPServer) Post(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPost, path, handleFunc)
}

func (s *HTTPServer) Get(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodGet, path, handleFunc)
}

// UseV1 会执行路由匹配，只有匹配上了的 mdls 才会生效
// 这个只需要稍微改造一下路由树就可以实现
func (s *HTTPServer) UseV1(path string, mdls ...Middleware) {
	panic("implement me")
}

func (s *HTTPServer) Use(mdls ...Middleware) {
	if s.mdls == nil {
		s.mdls = mdls
		return
	}
	s.mdls = append(s.mdls, mdls...)
}

func (s *HTTPServer) Response(ctx *Context) {
	if ctx.RespStatusCode > zero {
		ctx.Response.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Response.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败", err)
	}
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Request.Method, ctx.Request.URL.Path)
	if !ok || mi.n.handler == nil {
		// 没找到路由树 or 路由树未定义方法
		ctx.RespStatusCode = notFound
		return
	}
	ctx.PathParams = mi.pathParams
	// 命中的路由需要缓存起来
	ctx.MatchedRoute = mi.n.route
	mi.n.handler(ctx)
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:  req,
		Response: writer,
	}

	// 最后一个应该是 HTTPServer 执行路由匹配，执行用户代码
	root := s.serve
	// 从后往前组装
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}

	// 第一个应该是回写响应的
	// 因为它在调用next之后才回写响应，
	// 所以实际上 Response 是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.Response(ctx)
		}
	}
	root = m(root)
	root(ctx)

	// s.serve(ctx)
}

func (s *HTTPServer) Start(addr string) error {

	linstener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	println("成功监听端口")

	return http.Serve(linstener, s)

	// return http.ListenAndServe(addr, s)

}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

type HTTPServer struct {
	router
	mdls []Middleware
}

var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error
	addRoute(method, path string, handleFunc HandleFunc)
}
