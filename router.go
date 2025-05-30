package grog

import "github.com/startracex/grog/tire"

// HandlersNest pattern -> method -> []HandlerFunc
type HandlersNest map[string]map[string][]HandlerFunc

func (h HandlersNest) append(pattern, method string, handlers []HandlerFunc) {
	if _, ok := h[pattern]; !ok {
		h[pattern] = make(map[string][]HandlerFunc)
	}
	h[pattern][method] = append(h[pattern][method], handlers...)
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
	root     *tire.RouteTree
	handlers HandlersNest
}

// NewRouter create empty Router
func NewRouter() *Router {
	return &Router{
		root:     &tire.RouteTree{},
		handlers: make(HandlersNest),
	}
}

// register add a pattern -> method -> handlers
func (r *Router) register(method string, pattern string, handlers []HandlerFunc) {
	r.root.Insert(pattern, tire.SplitSlash(pattern), 0)
	r.handlers.append(pattern, method, handlers)
}

// Handle request or not found
func (r *Router) Handle(req *InnerRequest, res *InnerResponse) {
	path := req.Path
	node := r.root.Search(tire.SplitSlash(path), 0)
	if node != nil {
		pattern := node.Pattern
		method := req.Method
		req.Pattern = pattern
		req.Params = tire.ParseParams(path, pattern)
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
