package goup

import (
	"github.com/startracex/goup/core"
	"regexp"
	"strings"
)

// HandlersNest pattern -> method -> []HandlerFunc
type HandlersNest map[string]map[string][]HandlerFunc

// push handlers
func (h HandlersNest) push(pattern, method string, handlers []HandlerFunc) {
	if _, ok := h[pattern]; !ok {
		h[pattern] = make(map[string][]HandlerFunc)
	}
	h[pattern][method] = append(h[pattern][method], handlers...)
}

// get handlers and their exist
func (h HandlersNest) get(pattern, method string) (bool, bool, []HandlerFunc) {
	if _, ok := h[pattern]; !ok {
		return false, false, nil
	}
	if _, ok := h[pattern][method]; !ok {
		return true, false, nil
	}
	return true, true, h[pattern][method]
}

// Router type
type Router struct {
	// tire node
	root     *core.Node
	handlers HandlersNest
	re       bool
}

// NewRouter create empty Router
func NewRouter() *Router {
	return &Router{
		root:     &core.Node{},
		handlers: make(HandlersNest),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	var parts []string
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
	vs := strings.Split(s, "/")
	var parts []string
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
		}
	}
	return parts
}

// AddRoute add a pattern -> method -> handlers
func (r *Router) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	if !r.re {
		if strings.Contains(pattern, "/:") || strings.Contains(pattern, "/*") {
			parts := SplitPattern(pattern)
			r.root.Insert(pattern, parts, 0)
		}
	}
	r.handlers.push(pattern, method, handlers)
}

// GetRoute get match dynamic node and params
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
	var key = req.Path
	method := req.Method
	node, params := r.GetRoute(key)

	if node != nil {
		// dynamic router
		key = node.Pattern
	}

	havePattern, haveMethod, handlers := r.handlers.get(key, method)
	if havePattern {
		if haveMethod {
			req.Params = params
			req.appendHandlers(handlers)
		} else {
			req.appendHandlers(req.Engine.noMethodHandler)
		}
	} else {
		req.appendHandlers(req.Engine.noRouteHandler)
	}

	req.Next(res)
}

func (r *Router) HandlePrefix(req *HttpRequest, res *HttpResponse, prefix string) {
	key := req.Path
	method := req.Method
	for regex, maps := range r.handlers {
		if len(regex) < len(prefix) || len(regex) < len(key) {
			continue
		}
		if match, _ := regexp.MatchString(regex[len(prefix):], key[len(prefix):]); match {
			handlers, ok := maps[method]
			if ok {
				req.appendHandlers(handlers)
			} else {
				req.appendHandlers(req.Engine.noMethodHandler)
			}
			req.Next(res)
			return
		}
	}
	req.appendHandlers(req.Engine.noRouteHandler)
	req.Next(res)
}
