package v5

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/:id",
		},
	}

	mockHandler := func(ctx *Context) {}
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"a": {path: "a", children: map[string]*node{
						"b": {path: "b", starChild: &node{
							path: "*", handler: mockHandler},
						},
					}},
				},
				handler: mockHandler},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

}

func (r router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}
	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}

	if len(n.children) == 0 {
		return "", true
	}

	if n.starChild != nil {
		str, ok := n.starChild.equal(y.starChild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.path, str), false
		}
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true

}
