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

// GET defines the method to add GET request
func (group *RoutesGroup[T]) GET(pattern string, handlers ...T) {
	group.AddRoute(GET, pattern, handlers)
}

// POST defines the method to add POST request
func (group *RoutesGroup[T]) POST(pattern string, handlers ...T) {
	group.AddRoute(POST, pattern, handlers)
}

// PUT defines the method to add PUT request
func (group *RoutesGroup[T]) PUT(pattern string, handlers ...T) {
	group.AddRoute(PUT, pattern, handlers)
}

// DELETE defines the method to add DELETE request
func (group *RoutesGroup[T]) DELETE(pattern string, handlers ...T) {
	group.AddRoute(DELETE, pattern, handlers)
}

// PATCH defines the method to add PATCH request
func (group *RoutesGroup[T]) PATCH(pattern string, handlers ...T) {
	group.AddRoute(PATCH, pattern, handlers)
}

// OPTIONS defines the method to add OPTIONS request
func (group *RoutesGroup[T]) OPTIONS(pattern string, handlers ...T) {
	group.AddRoute(OPTIONS, pattern, handlers)
}

// HEAD defines the method to add HEAD request
func (group *RoutesGroup[T]) HEAD(pattern string, handlers ...T) {
	group.AddRoute(HEAD, pattern, handlers)
}

// CONNECT defines the method to add CONNECT request
func (group *RoutesGroup[T]) CONNECT(pattern string, handlers ...T) {
	group.AddRoute(CONNECT, pattern, handlers)
}

// TRACE defines the method to add TRACE request
func (group *RoutesGroup[T]) TRACE(pattern string, handlers ...T) {
	group.AddRoute(TRACE, pattern, handlers)
}

// METHOD defines the method to add request
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

// ANY is alias of ALL
func (group *RoutesGroup[T]) ANY(pattern string, handlers ...T) {
	group.ALL(pattern, handlers...)
}

// NoRoute accept not handlers for not found route
func (group *RoutesGroup[T]) NoRoute(handlers ...T) {
	group.Engine.noRoute = handlers
}

// NoMethod accept not handlers for not found method
func (group *RoutesGroup[T]) NoMethod(handlers ...T) {
	group.Engine.noMethod = handlers
}
