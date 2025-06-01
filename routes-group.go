package grog

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

func (group *RoutesGroup[T]) AddRoute(method string, pattern string, handlers []T) {
	group.Engine.Routes.AddRoute(method, group.Prefix+pattern, handlers)
}

func (group *RoutesGroup[T]) Use(middlewares ...T) {
	group.Middlewares = append(group.Middlewares, middlewares...)
}
