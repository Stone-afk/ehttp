package v4

import (
	"fmt"
	"strings"
)

// child 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否是命中参数路由
// 第三个返回值 bool 代表是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	root, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return root, false, ok

}

// findRoute 查找对应的节点
// 注意，返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	segs := strings.Split(path[1:], "/")
	mi := &matchInfo{}
	for _, s := range segs {
		var matchParam bool
		if s == "" {
			return nil, false
		}
		root, matchParam, ok = root.childOf(s)
		if !ok {
			return nil, false
		}
		if matchParam {
			mi.addValue(root.path[1:], s)
		}
	}
	mi.n = root
	return mi, true
}

// childOrCreate 查找子节点，如果子节点不存在就创建一个
// 并且将子节点放回去了 children 中
func (n *node) childOrCreate(path string) *node {
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.starChild == nil {
			n.starChild = &node{path: path}
		}
		return n.starChild
	}

	// 以 : 开头，我们认为是参数路由
	if path[0] == ':' {
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.paramChild != nil {
			if n.paramChild.path != path {
				panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
			}
		} else {
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

// addRoute 注册路由。
// method 是 HTTP 方法
// path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
func (r *router) addRoute(method, path string, handleFunc HandleFunc) {
	if path == "" {
		panic("web: 路由是空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}
	root, ok := r.trees[method]
	if !ok {
		// 这是一个全新的 HTTP 方法，创建根节点
		root = &node{path: "/"}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handleFunc
		return
	}
	// 开始一段段处理
	segs := strings.Split(path[1:], "/")
	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(s)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	root.handler = handleFunc
}

func newRouter() router {
	return router{
		trees: make(map[string]*node, 12),
	}
}

type router struct {
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
	trees map[string]*node
}

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 通配符匹配
// 这是不回溯匹配
type node struct {
	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符 * 表达的节点，任意匹配
	starChild *node

	// 参数节点
	paramChild *node
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key, value string) {
	if m.pathParams != nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = make(map[string]string)
	}
	m.pathParams[key] = value
}
