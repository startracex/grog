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
func (h HandlersNest) get(pattern, method string) (bool, []HandlerFunc) {
	hp, ok := h[pattern]
	if ok {
		hpm, ok := hp[method]
		if ok {
			return true, hpm
		}
	}
	return false, nil
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

// AddRoute add a pattern -> method -> handlers
func (r *Router) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	if !r.re {
		r.root.Set(pattern)
	}
	r.handlers.push(pattern, method, handlers)
}

// GetRoute get match dynamic node and params
func (r *Router) GetRoute(path string) (*core.Node, map[string]string) {
	searchParts := core.SplitSlash(path)
	n := r.root.Get(path)
	if n != nil {
		params := make(map[string]string)
		parts := core.SplitPattern(n.Pattern)
		for index, part := range parts {
			isWild, wildKey, isMulti := core.WildOf(part)
			if isWild {
				if isMulti {
					params[wildKey] = strings.Join(searchParts[index:], "/")
					break
				} else {
					params[wildKey] = searchParts[index]
				}
			}
		}
		return n, params
	}

	return nil, nil
}

// Handle request or not found
func (r *Router) Handle(req *HttpRequest, res *HttpResponse) {
	node, params := r.GetRoute(req.Path)
	if node != nil {
		pattern := node.Pattern
		method := req.Method

		haveMethod, handlers := r.handlers.get(pattern, method)
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
