package grog

import (
	"os"
	"path"
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
func (group *RoutesGroup) GET(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(GET, pattern, handlers)
}

// POST defines the method to add POST request
func (group *RoutesGroup) POST(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(POST, pattern, handlers)
}

// PUT defines the method to add PUT request
func (group *RoutesGroup) PUT(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(PUT, pattern, handlers)
}

// DELETE defines the method to add DELETE request
func (group *RoutesGroup) DELETE(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(DELETE, pattern, handlers)
}

// PATCH defines the method to add PATCH request
func (group *RoutesGroup) PATCH(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(PATCH, pattern, handlers)
}

// OPTIONS defines the method to add OPTIONS request
func (group *RoutesGroup) OPTIONS(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(OPTIONS, pattern, handlers)
}

// HEAD defines the method to add HEAD request
func (group *RoutesGroup) HEAD(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(HEAD, pattern, handlers)
}

// CONNECT defines the method to add CONNECT request
func (group *RoutesGroup) CONNECT(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(CONNECT, pattern, handlers)
}

// TRACE defines the method to add TRACE request
func (group *RoutesGroup) TRACE(pattern string, handlers ...HandlerFunc) {
	group.AddRoute(TRACE, pattern, handlers)
}

// METHOD defines the method to add request
func (group *RoutesGroup) METHOD(method, pattern string, handlers ...HandlerFunc) {
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
func (group *RoutesGroup) ALL(pattern string, handlers ...HandlerFunc) {
	for _, method := range AllMethods {
		group.AddRoute(method, pattern, handlers)
	}
}

// ANY is alias of ALL
func (group *RoutesGroup) ANY(pattern string, handlers ...HandlerFunc) {
	group.ALL(pattern, handlers...)
}

// NoRoute accept not handlers for not found route
func (group *RoutesGroup) NoRoute(handlers ...HandlerFunc) {
	group.Engine.NoRouteHandler = handlers
}

// NoMethod accept not handlers for not found method
func (group *RoutesGroup) NoMethod(handlers ...HandlerFunc) {
	group.Engine.NoMethodHandler = handlers
}

// // Public handle file, or directory
func (group *RoutesGroup) Public(pattern string, path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if fi.IsDir() {
		group.Directory(pattern, path)
	} else {
		group.File(pattern, path)
	}
}

// Static is alias of Public
func (group *RoutesGroup) Static(pattern string, path string) {
	group.Public(pattern, path)
}

// Directory handle directory
func (group *RoutesGroup) Directory(pattern string, root string) {
	key := "path"
	handler := func(c *Context) {
		p := c.Params[key]
		p = strings.ReplaceAll(p, "\\", "/")
		p = path.Clean(p)
		p = strings.TrimPrefix(p, "/")
		ServeFile(c, path.Join(root, p))
	}
	pattern = path.Join(pattern, "*"+key)
	group.GET(pattern, handler)
	group.HEAD(pattern, handler)
}

// File handle file
func (group *RoutesGroup) File(pattern string, filepath string) {
	handler := func(c *Context) {
		ServeFile(c, filepath)
	}
	group.GET(pattern, handler)
	group.HEAD(pattern, handler)
}
