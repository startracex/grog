package grog

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/startracex/grog/dns"
)

var Host = "127.0.0.1"

type Engine[T any] struct {
	*RoutesGroup[T]
	Routes      *Routes[T]
	groups      []*RoutesGroup[T]
	noRoute     []T
	noMethod    []T
	DNS         *dns.DNS[*Engine[T]]
	Adapter     func(T) func(Context)
	ContextPool sync.Pool
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

	var c *handleContext[T]
	if v := e.ContextPool.Get(); v != nil {
		c = v.(*handleContext[T])
	} else {
		c = new(handleContext[T])
	}
	c.request = req
	c.writer = res
	c.index = -1
	c.adapter = e.Adapter
	path := req.URL.Path

	node := e.Routes.Search(path)
	if node != nil {
		c.node = node
		c.pattern = node.Pattern
		handlers, ok := node.Value[req.Method]
		if !ok {
			c.handlers = e.noMethod
		} else {
			for _, group := range e.groups {
				if strings.HasPrefix(c.pattern, group.Prefix+"/") {
					c.handlers = append(c.handlers, group.Middlewares...)
				}
			}
			c.handlers = append(c.handlers, handlers...)
		}
	} else {
		c.handlers = e.noMethod
	}

	c.Next()

	c.request = nil
	c.writer = nil
	c.pattern = ""
	c.index = -1
	c.handlers = c.handlers[:0]
	c.node = nil
	e.ContextPool.Put(c)
}

// New create engine
func New[T any]() *Engine[T] {
	engine := &Engine[T]{
		Routes: NewRouter[T](),
		ContextPool: sync.Pool{
			New: func() any {
				return new(handleContext[T])
			},
		},
	}
	engine.RoutesGroup = &RoutesGroup[T]{Engine: engine}
	engine.groups = []*RoutesGroup[T]{engine.RoutesGroup}
	engine.Adapter = defaultAdapter
	return engine
}

func (e *Engine[T]) NoMethod(handlers ...T) []T {
	e.noMethod = append(e.noMethod, handlers...)
	return e.noMethod
}

func (e *Engine[T]) NoRoute(handlers ...T) []T {
	e.noRoute = append(e.noRoute, handlers...)
	return e.noRoute
}

func (e *Engine[T]) Domain(domains ...string) *Engine[T] {
	newEngine := New[T]()
	newEngine.noMethod = e.noMethod
	newEngine.noRoute = e.noRoute
	newEngine.Use(e.Middlewares...)

	if e.DNS == nil {
		e.DNS = dns.NewDNS[*Engine[T]]()
	}

	for _, domain := range domains {
		e.DNS.Insert(domain, newEngine)
	}
	return newEngine
}

func normalizeAddr(addr any) string {
	return ":" + strings.Trim(fmt.Sprintf("%v", addr), ":")
}

// Run call ListenAndServe
func (e *Engine[T]) Run(addr any) error {
	return e.ListenAndServe(addr)
}

// RunTLS call ListenAndServeTLS
func (e *Engine[T]) RunTLS(addr any, cert, key string) error {
	return e.ListenAndServeTLS(addr, cert, key)
}

// ListenAndServe start a server
func (e *Engine[T]) ListenAndServe(addr any) error {
	return http.ListenAndServe(normalizeAddr(addr), e)
}

// ListenAndServeTLS start a server with TLS
func (e *Engine[T]) ListenAndServeTLS(addr any, cert, key string) error {
	return http.ListenAndServeTLS(normalizeAddr(addr), cert, key, e)
}
