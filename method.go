package goup

import (
    "net/http"
    "os"
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
    ANY     = "ANY"
)

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(GET, pattern, handlers)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(POST, pattern, handlers)
}

// PUT defines the method to add PUT request
func (group *RouterGroup) PUT(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(PUT, pattern, handlers)
}

// DELETE defines the method to add DELETE request
func (group *RouterGroup) DELETE(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(DELETE, pattern, handlers)
}

// PATCH defines the method to add PATCH request
func (group *RouterGroup) PATCH(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(PATCH, pattern, handlers)
}

// OPTIONS defines the method to add OPTIONS request
func (group *RouterGroup) OPTIONS(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(OPTIONS, pattern, handlers)
}

// HEAD defines the method to add HEAD request
func (group *RouterGroup) HEAD(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(HEAD, pattern, handlers)
}

// CONNECT defines the method to add CONNECT request
func (group *RouterGroup) CONNECT(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(CONNECT, pattern, handlers)
}

// TRACE defines the method to add TRACE request
func (group *RouterGroup) TRACE(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(TRACE, pattern, handlers)
}

// METHOD defines the method to add request
func (group *RouterGroup) METHOD(method, pattern string, handlers ...HandlerFunc) {
    group.AddRoute(strings.ToUpper(method), pattern, handlers)
}

// ALL defines the method to add all requests
func (group *RouterGroup) ALL(pattern string, handlers ...HandlerFunc) {
    all := []string{GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, CONNECT, TRACE}
    for _, method := range all {
        group.AddRoute(method, pattern, handlers)
    }
}

// ANY defines the method to add any request
func (group *RouterGroup) ANY(pattern string, handlers ...HandlerFunc) {
    group.AddRoute(ANY, pattern, handlers)
}

// NoRoute accept not found subpath handler
func (group *RouterGroup) NoRoute(handlers ...HandlerFunc) {
    group.ANY("/*url", handlers...)
}

// File handle file/directory
func (group *RouterGroup) File(pattern string, path string) {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    fi, _ := f.Stat()
    if err != nil {
        panic(err)
    }
    if fi.IsDir() {
        group.Static(pattern, path)
    } else {
        group.GET(pattern, func(req *HttpRequest, res *HttpResponse) {
            http.ServeFile(res.Writer, req.OriginalRequest, path)
        })
        group.HEAD(pattern, func(req *HttpRequest, res *HttpResponse) {
            http.ServeFile(res.Writer, req.OriginalRequest, path)
        })
    }
}

// Static handle directory
func (group *RouterGroup) Static(pattern string, root string) {
    key := "path"
    handler := func(req *HttpRequest, res *HttpResponse) {
        file := req.Params[key]
        if len(file) == 0 {
            res.Status(404)
            return
        }
        http.ServeFile(res.Writer, req.OriginalRequest, root+"/"+file)
    }
    group.GET(pattern+"/*"+key, handler)
    group.HEAD(pattern+"/*"+key, handler)
}
