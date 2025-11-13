package grog

import (
	"errors"

	"github.com/startracex/grog/router"
)

type Routes[T any] struct {
	Root *router.Router[map[string][]T]
}

func NewRoutes[T any]() *Routes[T] {
	return &Routes[T]{
		Root: router.New[map[string][]T](),
	}
}

func (r *Routes[T]) AddRoute(method string, pattern string, handlers []T) {
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
	r.Root.Insert(pattern, map[string][]T{method: handlers})
}

var ErrNoRoute = errors.New("grog: no route")
var ErrNoMethod = errors.New("grog: no method")

func (r *Routes[T]) Search(path string) *router.Router[map[string][]T] {
	return r.Root.Search(path)
}

type RoutesGroup[T any] struct {
	Prefix      string
	Middlewares []T
	Engine      *Engine[T]
}

func (group *RoutesGroup[T]) Group(prefix string, middlewares ...T) *RoutesGroup[T] {
	engine := group.Engine
	newGroup := &RoutesGroup[T]{
		Prefix:      group.Prefix + prefix,
		Engine:      engine,
		Middlewares: middlewares,
	}
	newGroup.Middlewares = append(newGroup.Middlewares, middlewares...)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RoutesGroup[T]) AddRoute(method string, pattern string, handlers []T) *RoutesGroup[T] {
	group.Engine.Routes.AddRoute(method, group.Prefix+pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) Use(middlewares ...T) *RoutesGroup[T] {
	group.Middlewares = append(group.Middlewares, middlewares...)
	return group
}
