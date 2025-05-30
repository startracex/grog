package grog

type RouterGroup struct {
	Prefix      string
	Middlewares []HandlerFunc
	Engine *Engine
}

func (group *RouterGroup) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup {
	engine := group.Engine
	newGroup := &RouterGroup{
		Prefix: group.Prefix + prefix,
		Engine:      engine,
		Middlewares: middlewares,
	}
	newGroup.Middlewares = append(newGroup.Middlewares, middlewares...)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	group.Engine.Routes.AddRoute(method, group.Prefix+pattern, handlers)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Middlewares = append(group.Middlewares, middlewares...)
}
