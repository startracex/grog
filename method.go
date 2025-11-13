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

func (group *RoutesGroup[T]) GET(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(GET, pattern, handlers)
}

func (group *RoutesGroup[T]) POST(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(POST, pattern, handlers)
}

func (group *RoutesGroup[T]) PUT(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(PUT, pattern, handlers)
}

func (group *RoutesGroup[T]) DELETE(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(DELETE, pattern, handlers)
}

func (group *RoutesGroup[T]) PATCH(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(PATCH, pattern, handlers)
}

func (group *RoutesGroup[T]) OPTIONS(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(OPTIONS, pattern, handlers)
}

func (group *RoutesGroup[T]) HEAD(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(HEAD, pattern, handlers)
}

func (group *RoutesGroup[T]) CONNECT(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(CONNECT, pattern, handlers)
}

func (group *RoutesGroup[T]) TRACE(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.AddRoute(TRACE, pattern, handlers)
}

var allMethods = [9]string{GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, CONNECT, TRACE}

func (group *RoutesGroup[T]) METHOD(method, pattern string, handlers ...T) *RoutesGroup[T] {
	method = strings.ToUpper(method)
	for _, m := range allMethods {
		if method == m {
			return group.AddRoute(method, pattern, handlers)
		}
	}
	panic("unsupported method")
}

// ALL defines the method to add all requests
func (group *RoutesGroup[T]) ALL(pattern string, handlers ...T) *RoutesGroup[T] {
	for _, method := range allMethods {
		group.AddRoute(method, pattern, handlers)
	}
	return group
}

func (group *RoutesGroup[T]) ANY(pattern string, handlers ...T) *RoutesGroup[T] {
	return group.ALL(pattern, handlers...)
}
