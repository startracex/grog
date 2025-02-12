package goup

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

var Host = "127.0.0.1"

type Engine struct {
	*RouterGroup
	router          *Router
	groups          []*RouterGroup
	Pool            sync.Pool
	Template        template.Template
	FuncMap         template.FuncMap
	noRouteHandler  []HandlerFunc
	noMethodHandler []HandlerFunc
}

// ServeHTTP for http.ListenAndServe
func (e *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	newRequest := NewRequest(req)
	newRequest.Engine = e
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix+"/") {
			newRequest.appendHandlers(group.middlewares)

		}
	}

	newResponse := NewResponse(res)
	e.router.Handle(&newRequest, &newResponse)
}

// New create engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.Pool.New = func() any {
		return bytes.NewBuffer(make([]byte, 4096))
	}
	engine.NoRoute(func(request Request, response Response) {
		response.StatusText(404)
	})
	engine.NoMethod(func(request Request, response Response) {
		response.StatusText(405)
	})
	return engine
}

// Default use DefaultMiddleware
func Default() *Engine {
	engine := New()
	engine.Use(DefaultMiddleware...)
	return engine
}

// BaseURL set engine's prefix
func (e *Engine) BaseURL(base string) {
	if base != "/" {
		e.prefix = base
	}
}

// SetPoolNew Replace the default NEW func
func (e *Engine) SetPoolNew(f func() any) {
	e.Pool.New = f
}

// LoadHTMLFiles load the path file
func (e *Engine) LoadHTMLFiles(path ...string) {
	e.Template = *template.Must(template.New("").Funcs(e.FuncMap).ParseFiles(path...))
}

// LoadFunc load func map for template
func (e *Engine) LoadFunc(funcMap template.FuncMap) {
	e.FuncMap = funcMap
}

// ListenAndServe start a server
func (e *Engine) ListenAndServe(addr string) error {
	return http.ListenAndServe(mustPort(addr), e)
}

// Run call ListenAndServe or ListenAndServeTLS if it has filePath slice
func (e *Engine) Run(addr string, filePath ...string) error {
	addr = mustPort(addr)
	if len(filePath) > 1 {
		fmt.Println("Listen and serve TLS at https://" + Host + addr)
		return e.ListenAndServeTLS(addr, filePath[0], filePath[1])
	}
	fmt.Println("Listen and serve at http://" + Host + addr)
	return e.ListenAndServe(addr)
}

// ListenAndServeTLS start a server with TLS
func (e *Engine) ListenAndServeTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(mustPort(addr), cert, key, e)
}

// mustPort make sure addr is a valid port
func mustPort(addr string) string {
	if !strings.HasPrefix(addr, ":") {
		return ":" + addr
	}
	return addr
}
