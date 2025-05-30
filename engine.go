package grog

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/startracex/grog/dns"
	"github.com/startracex/grog/router"
)

var Host = "127.0.0.1"

type Engine struct {
	*RouterGroup
	Routes          *Routes
	groups          []*RouterGroup
	Pool            sync.Pool
	Template        *template.Template
	NoRouteHandler  []HandlerFunc
	NoMethodHandler []HandlerFunc
	DNS             *dns.DNS[*Engine]
}

// ServeHTTP for http.ListenAndServe
func (e *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if e.DNS != nil {
		domain := dns.GetDomain(req.Host)
		matchEngine, ok := e.DNS.Match(domain)
		if ok {
			matchEngine.ServeHTTP(res, req)
			return
		}
	}

	path := req.URL.Path
	var pattern string
	var params map[string]string
	var methods []string
	hf := make([]HandlerFunc, 0)

	node := e.Routes.Search(path)
	if node != nil {
		pattern = node.Pattern
		handler, ok := node.Value[req.Method]
		if !ok {
			hf = e.NoMethodHandler
		} else {
			methods = make([]string, len(node.Value))
			for k := range node.Value {
				methods = append(methods, k)
			}
			for _, group := range e.groups {
				if strings.HasPrefix(path, group.Prefix+"/") {
					hf = append(hf, group.Middlewares...)
				}
			}
			hf = append(hf, handler...)
		}
		params = router.ParseParams(path, node.Pattern)
	} else {
		hf = e.NoMethodHandler
	}

	c := &Context{
		Request:  req,
		Writer:   res,
		Engine:   e,
		Handlers: hf,
		Index:    -1,
		Pattern:  pattern,
		Params:   params,
		Methods:  methods,
	}
	c.Next()
}

// New create engine
func New() *Engine {
	engine := &Engine{
		Routes: NewRouter(),
		Pool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 4096))
			},
		},
	}
	engine.RouterGroup = &RouterGroup{Engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) Domain(domains ...string) *Engine {
	newEngine := New()
	newEngine.NoMethodHandler = e.NoMethodHandler
	newEngine.NoRouteHandler = e.NoRouteHandler
	newEngine.Use(e.Middlewares...)

	if e.DNS == nil {
		e.DNS = dns.NewDNS[*Engine]()
	}

	for _, domain := range domains {
		e.DNS.Insert(domain, newEngine)
	}
	return newEngine
}

// PoolNew Replace the default NEW func
func (e *Engine) PoolNew(f func() any) {
	e.Pool.New = f
}

// ParseTemplateFiles load the path file
func (e *Engine) ParseTemplateFiles(funcMap template.FuncMap, path ...string) {
	if len(path) == 0 {
		return
	}
	e.Template = template.Must(template.New("").Funcs(funcMap).ParseFiles(path...))
}

func normalizeAddr(addr string) string {
	if !strings.HasPrefix(addr, ":") {
		return ":" + addr
	}
	return addr
}

// ListenAndServe start a server
func (e *Engine) ListenAndServe(addr string) error {
	return http.ListenAndServe(normalizeAddr(addr), e)
}

// Run call ListenAndServe
func (e *Engine) Run(addr string) error {
	return e.ListenAndServe(addr)
}

// RunTLS call ListenAndServeTLS
func (e *Engine) RunTLS(addr, cert, key string) error {
	return e.ListenAndServeTLS(addr, cert, key)
}

// ListenAndServeTLS start a server with TLS
func (e *Engine) ListenAndServeTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(normalizeAddr(addr), cert, key, e)
}
