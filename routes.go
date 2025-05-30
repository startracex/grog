package grog

import (
	"errors"

	"github.com/startracex/grog/router"
)

type Routes struct {
	Root *router.Router[map[string][]HandlerFunc]
}

func NewRouter() *Routes {
	return &Routes{
		Root: router.NewRouter[map[string][]HandlerFunc](),
	}
}

func (r *Routes) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	node := r.Root.SearchPattern(pattern)
	if node != nil {
		m := node.Value
		if _, ok := m[method]; ok {
			m[method] = append(m[method], handlers...)
		} else {
			m[method] = handlers
		}
		return
	}
	r.Root.InsertPattern(pattern, map[string][]HandlerFunc{
		method: handlers,
	})
}

var ErrNoRoute = errors.New("grog: no route")
var ErrNoMethod = errors.New("grog: no method")

func (r *Routes) GetHandlers(pattern, method string) ([]HandlerFunc, error) {
	node := r.Root.SearchPattern(pattern)
	if node != nil {
		handlers, ok := node.Value[method]
		if ok {
			return handlers, nil
		}
		return nil, ErrNoMethod
	}
	return nil, ErrNoRoute
}

func (r *Routes) Handle(req *InnerRequest, res *InnerResponse) {
	path := req.Path
	node := r.Root.SearchPattern(path)
	if node != nil {
		pattern := node.Pattern
		method := req.Method
		req.Pattern = pattern
		req.Params = router.ParseParams(path, pattern)
		handlers, ok := node.Value[method]
		if ok {
			req.appendHandlers(handlers)
		} else {
			req.appendHandlers(req.Engine.NoMethodHandler)
		}
	} else {
		req.appendHandlers(req.Engine.NoRouteHandler)
	}

	req.Next(res)
}

func (r *Routes) Has(pattern, method string) (bool, bool) {
	node := r.Root.SearchPattern(pattern)
	if node != nil {
		_, ok := node.Value[method]
		return true, ok
	}
	return false, false
}

func (r *Routes) AllMethods(pattern string) []string {
	node := r.Root.SearchPattern(pattern)
	if node != nil {
		s := make([]string, len(node.Value))
		for k := range node.Value {
			s = append(s, k)
		}
		return s
	}
	return nil
}
