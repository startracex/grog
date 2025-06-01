package grog

import (
	"bytes"
	"net/http"
	"strings"
	"sync"

	"github.com/startracex/grog/dns"
	"github.com/startracex/grog/router"
)

var Host = "127.0.0.1"

type Engine[T any] struct {
	*RoutesGroup[T]
	Routes          *Routes[T]
	groups          []*RoutesGroup[T]
	Pool            sync.Pool
	NoRouteHandler  []T
	NoMethodHandler []T
	DNS             *dns.DNS[*Engine[T]]
	Adapter         func(T) func(Context)
}

// ServeHTTP for http.ListenAndServe
func (e *Engine[T]) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
	hf := make([]T, 0)

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

	c := &HandleContext[T]{
		request:        req,
		writer:         res,
		HandlerAdapter: e.Adapter,
		Handlers:       hf,
		Index:          -1,
		Pattern:        pattern,
		Params:         params,
		Methods:        methods,
	}
	c.Next()
}

// New create engine
func New[T any]() *Engine[T] {
	engine := &Engine[T]{
		Routes: NewRouter[T](),
		Pool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 4096))
			},
		},
	}
	engine.RoutesGroup = &RoutesGroup[T]{Engine: engine}
	engine.groups = []*RoutesGroup[T]{engine.RoutesGroup}
	return engine
}

func (e *Engine[T]) Domain(domains ...string) *Engine[T] {
	newEngine := New[T]()
	newEngine.NoMethodHandler = e.NoMethodHandler
	newEngine.NoRouteHandler = e.NoRouteHandler
	newEngine.Use(e.Middlewares...)

	if e.DNS == nil {
		e.DNS = dns.NewDNS[*Engine[T]]()
	}

	for _, domain := range domains {
		e.DNS.Insert(domain, newEngine)
	}
	return newEngine
}

// PoolNew Replace the default NEW func
func (e *Engine[T]) PoolNew(f func() any) {
	e.Pool.New = f
}

func normalizeAddr(addr string) string {
	if !strings.HasPrefix(addr, ":") {
		return ":" + addr
	}
	return addr
}

// ListenAndServe start a server
func (e *Engine[T]) ListenAndServe(addr string) error {
	return http.ListenAndServe(normalizeAddr(addr), e)
}

// Run call ListenAndServe
func (e *Engine[T]) Run(addr string) error {
	return e.ListenAndServe(addr)
}

// RunTLS call ListenAndServeTLS
func (e *Engine[T]) RunTLS(addr, cert, key string) error {
	return e.ListenAndServeTLS(addr, cert, key)
}

// ListenAndServeTLS start a server with TLS
func (e *Engine[T]) ListenAndServeTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(normalizeAddr(addr), cert, key, e)
}
