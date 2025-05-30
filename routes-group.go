package grog

type RoutesGroup struct {
	Prefix      string
	Middlewares []HandlerFunc
	Engine      *Engine
}

func (group *RoutesGroup) Group(prefix string, middlewares ...HandlerFunc) *RoutesGroup {
	engine := group.Engine
	newGroup := &RoutesGroup{
		Prefix:      group.Prefix + prefix,
		Engine:      engine,
		Middlewares: middlewares,
	}
	newGroup.Middlewares = append(newGroup.Middlewares, middlewares...)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RoutesGroup) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	group.Engine.Routes.AddRoute(method, group.Prefix+pattern, handlers)
}

func (group *RoutesGroup) Use(middlewares ...HandlerFunc) {
	group.Middlewares = append(group.Middlewares, middlewares...)
}
