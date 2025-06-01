package grog

import (
	"strings"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
	CONNECT = "CONNECT"
)

var AllMethods = [...]string{GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, CONNECT, TRACE}

func (group *RoutesGroup[T]) GET(pattern string, handlers ...T) {
	group.AddRoute(GET, pattern, handlers)
}

func (group *RoutesGroup[T]) POST(pattern string, handlers ...T) {
	group.AddRoute(POST, pattern, handlers)
}

func (group *RoutesGroup[T]) PUT(pattern string, handlers ...T) {
	group.AddRoute(PUT, pattern, handlers)
}

func (group *RoutesGroup[T]) DELETE(pattern string, handlers ...T) {
	group.AddRoute(DELETE, pattern, handlers)
}

func (group *RoutesGroup[T]) PATCH(pattern string, handlers ...T) {
	group.AddRoute(PATCH, pattern, handlers)
}

func (group *RoutesGroup[T]) OPTIONS(pattern string, handlers ...T) {
	group.AddRoute(OPTIONS, pattern, handlers)
}

func (group *RoutesGroup[T]) HEAD(pattern string, handlers ...T) {
	group.AddRoute(HEAD, pattern, handlers)
}

func (group *RoutesGroup[T]) CONNECT(pattern string, handlers ...T) {
	group.AddRoute(CONNECT, pattern, handlers)
}

func (group *RoutesGroup[T]) TRACE(pattern string, handlers ...T) {
	group.AddRoute(TRACE, pattern, handlers)
}

func (group *RoutesGroup[T]) METHOD(method, pattern string, handlers ...T) {
	method = strings.ToUpper(method)
	for _, m := range AllMethods {
		if method == m {
			group.AddRoute(method, pattern, handlers)
			return
		}
	}
	panic("unsupported method")
}

// ALL defines the method to add all requests
func (group *RoutesGroup[T]) ALL(pattern string, handlers ...T) {
	for _, method := range AllMethods {
		group.AddRoute(method, pattern, handlers)
	}
}

func (group *RoutesGroup[T]) ANY(pattern string, handlers ...T) {
	group.ALL(pattern, handlers...)
}

func (group *RoutesGroup[T]) NoRoute(handlers ...T) {
	group.Engine.noRoute = handlers
}

func (group *RoutesGroup[T]) NoMethod(handlers ...T) {
	group.Engine.noMethod = handlers
}
