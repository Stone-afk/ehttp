package v2

import (
	"net"
	"net/http"
)

func (s *HTTPServer) Post(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPost, path, handleFunc)
}

func (s *HTTPServer) Get(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodGet, path, handleFunc)
}

func (s *HTTPServer) serve(ctx *Context) {
	node, ok := s.findRoute(ctx.Request.Method, ctx.Request.URL.Path)
	if !ok || node.handler != nil {
		ctx.Response.WriteHeader(404)
		ctx.Response.Write([]byte("Not Found"))
		return
	}
	node.handler(ctx)
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:  req,
		Response: writer,
	}

	s.serve(ctx)
}

func (s *HTTPServer) Start(addr string) error {

	linstener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	return http.Serve(linstener, s)

}

type HTTPServer struct {
	router
}

var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error
	addRoute(method, path string, handleFunc HandleFunc)
}
