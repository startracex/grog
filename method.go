package grog

import (
	"net/http"
	"strings"
)

const (
	GET     = http.MethodGet
	POST    = http.MethodPost
	PUT     = http.MethodPut
	DELETE  = http.MethodDelete
	PATCH   = http.MethodPatch
	HEAD    = http.MethodHead
	OPTIONS = http.MethodOptions
	TRACE   = http.MethodTrace
	CONNECT = http.MethodConnect
)

var AllMethods = [...]string{GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, CONNECT, TRACE}

func (group *RoutesGroup[T]) GET(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(GET, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) POST(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(POST, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) PUT(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(PUT, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) DELETE(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(DELETE, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) PATCH(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(PATCH, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) OPTIONS(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(OPTIONS, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) HEAD(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(HEAD, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) CONNECT(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(CONNECT, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) TRACE(pattern string, handlers ...T) *RoutesGroup[T] {
	group.AddRoute(TRACE, pattern, handlers)
	return group
}

func (group *RoutesGroup[T]) METHOD(method, pattern string, handlers ...T) *RoutesGroup[T] {
	method = strings.ToUpper(method)
	for _, m := range AllMethods {
		if method == m {
			group.AddRoute(method, pattern, handlers)
			return group
		}
	}
	panic("unsupported method")
}

// ALL defines the method to add all requests
func (group *RoutesGroup[T]) ALL(pattern string, handlers ...T) *RoutesGroup[T] {
	for _, method := range AllMethods {
		group.AddRoute(method, pattern, handlers)
	}
	return group
}

func (group *RoutesGroup[T]) ANY(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.ALL(pattern, handlers...)
}
