package goup

import (
	"github.com/startracex/goup/core"
	"strings"
)

type handlersInMethod = map[string][]HandlerFunc

type Router struct {
	root     *core.Node
	handlers map[string]handlersInMethod
}

func NewRouter() *Router {
	return &Router{
		root:     &core.Node{},
		handlers: make(map[string]handlersInMethod),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func SplitPattern(s string) []string {
	return parsePattern(s)
}

func SplitSlash(s string) []string {
	return strings.Split(s, "/")[1:]
}

func (r *Router) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	parts := SplitPattern(pattern)
	r.root.Insert(pattern, parts, 0)
	if r.handlers[pattern] == nil {
		r.handlers[pattern] = make(map[string][]HandlerFunc)
	}
	r.handlers[pattern][method] = handlers
}

func (r *Router) GetRoute(path string) (*core.Node, map[string]string) {
	searchParts := SplitSlash(path)
	n := r.root.Search(searchParts, 0)
	if n != nil {
		params := make(map[string]string)
		parts := SplitPattern(n.Pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// Handle request or not found
func (r *Router) Handle(req *HttpRequest, res *HttpResponse) {
	method := req.Method
	path := req.URL().Path
	node, params := r.GetRoute(path)
	if node != nil {
		handlers, ok := r.handlers[node.Pattern][method]
		if ok {
			req.Params = params
			req.Handlers = append(req.Handlers, handlers...)
		} else {
			req.Handlers = req.Engine.noMethodHandler
		}
	} else {
		req.Handlers = req.Engine.noRouteHandler
	}
	req.Next(res)
}
