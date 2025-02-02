package goup

import (
	"regexp"
)

// HandlersNest pattern -> method -> []HandlerFunc
type HandlersNest map[string]map[string][]HandlerFunc

// append handlers
func (h HandlersNest) append(pattern, method string, handlers []HandlerFunc) {
	if _, ok := h[pattern]; !ok {
		h[pattern] = make(map[string][]HandlerFunc)
	}
	h[pattern][method] = append(h[pattern][method], handlers...)
}

func (h HandlersNest) hasPattern(pattern string) bool {
	_, ok := h[pattern]
	return ok
}

func (h HandlersNest) hasMethod(pattern, method string) bool {
	if h.hasPattern(pattern) {
		_, ok := h[pattern][method]
		return ok
	}
	return false
}

func (h HandlersNest) allPatterns() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h HandlersNest) allMethods(pattern string) []string {
	hp, ok := h[pattern]
	if ok {
		keys := make([]string, 0, len(hp))
		for k := range hp {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

// Router type
type Router struct {
	root     *RouteTree
	handlers HandlersNest
	re       bool
}

// NewRouter create empty Router
func NewRouter() *Router {
	return &Router{
		root:     &RouteTree{},
		handlers: make(HandlersNest),
	}
}

// register add a pattern -> method -> handlers
func (r *Router) register(method string, pattern string, handlers []HandlerFunc) {
	if !r.re {
		r.root.Insert(pattern, SplitSlash(pattern), 0)
	}
	r.handlers.append(pattern, method, handlers)
}

// Handle request or not found
func (r *Router) Handle(req *InnerRequest, res *InnerResponse) {
	path := req.Path
	node := r.root.Search(SplitSlash(path), 0)
	if node != nil {
		pattern := node.Pattern
		method := req.Method
		req.Pattern = pattern
		req.Params = ParseParams(path, pattern)
		handlers, ok := r.handlers[pattern][method]
		if ok {
			req.appendHandlers(handlers)
		} else {
			req.appendHandlers(req.Engine.noMethodHandler)
		}
	} else {
		req.appendHandlers(req.Engine.noRouteHandler)
	}

	req.Next(res)
}

func (r *Router) HandlePrefix(req *InnerRequest, res *InnerResponse, prefix string) {
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
