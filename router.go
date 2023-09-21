package goup

import (
	"github.com/startracex/goup/core"
	"strings"
)

type Router struct {
	roots    map[string]*core.Node
	handlers map[string][]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*core.Node),
		handlers: make(map[string][]HandlerFunc),
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

func (r *Router) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &core.Node{}
	}
	r.roots[method].Insert(pattern, parts, 0)
	r.handlers[key] = handlers
}

func (r *Router) getRoute(method string, path string) (*core.Node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.Search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.Pattern)
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
/* 
func (r *Router) getRoutes(method string) []*core.Node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*core.Node, 0)
	root.Travel(&nodes)
	return nodes
}
 */
// Handle request or not found
func (r *Router) Handle(req *HttpRequest, res *HttpResponse) {
	method := req.Method
	path := req.URL().Path
	node, params := r.getRoute(method, path)
	if node == nil {
		node, params = r.getRoute(ANY, path)
		method = ANY
	}
	if node != nil {
		key := method + "-" + node.Pattern
		req.Params = params
		req.Handlers = append(req.Handlers, r.handlers[key]...)
	} else {
		res.Error(404, "NOT FOUND: "+path)
	}
	req.Next(res)
}
