package v1

import (
	"net"
	"net/http"
)

func (s *HTTPServer) Post(path string, handleFunc HandleFunc) {
	s.AddRoute(http.MethodPost, path, handleFunc)
}

func (s *HTTPServer) Get(path string, handleFunc HandleFunc) {
	s.AddRoute(http.MethodGet, path, handleFunc)
}

func (s *HTTPServer) AddRoute(method, path string, handleFunc HandleFunc) {
	panic("implement me")
}

func (s *HTTPServer) serve(ctx *Context) {

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
}

var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

//type Server interface {
//	http.Handler
//	//Start(addr string) error
//	//AddRoute(method, path string, handleFunc HandleFunc)
//}

type Server interface {
	http.Handler
}
