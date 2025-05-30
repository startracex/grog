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
	node := r.Search(pattern)
	if node != nil {
		m := node.Value
		if _, ok := m[method]; ok {
			m[method] = append(m[method], handlers...)
		} else {
			m[method] = handlers
		}
		return
	}
	r.Root.Insert(pattern, router.SplitSlash(pattern), 0, map[string][]HandlerFunc{
		method: handlers,
	})
}

var ErrNoRoute = errors.New("grog: no route")
var ErrNoMethod = errors.New("grog: no method")

func (r *Routes) Search(path string) *router.Router[map[string][]HandlerFunc] {
	parts := router.SplitSlash(path)
	return r.Root.Search(parts, 0)
}

func (r *Routes) AllMethods(pattern string) []string {
	node := r.Search(pattern)
	if node != nil {
		s := make([]string, len(node.Value))
		for k := range node.Value {
			s = append(s, k)
		}
		return s
	}
	return nil
}
