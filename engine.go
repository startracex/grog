package grog

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/startracex/grog/dns"
)

var Host = "127.0.0.1"

type Engine[T any] struct {
	*RoutesGroup[T]
	Routes   *Routes[T]
	groups   []*RoutesGroup[T]
	noRoute  []T
	noMethod []T
	DNS      *dns.DNS[*Engine[T]]
	Adapter  func(T) func(Context)
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
	var allowMethods []string
	hf := make([]T, 0)

	node := e.Routes.Search(path)
	if node != nil {
		pattern = node.Pattern
		handler, ok := node.Value[req.Method]
		if !ok {
			hf = e.noMethod
		} else {
			allowMethods = make([]string, len(node.Value))
			for k := range node.Value {
				allowMethods = append(allowMethods, k)
			}
			for _, group := range e.groups {
				if strings.HasPrefix(path, group.Prefix+"/") {
					hf = append(hf, group.Middlewares...)
				}
			}
			hf = append(hf, handler...)
		}
	} else {
		hf = e.noMethod
	}

	c := &HandleContext[T]{
		request:      req,
		writer:       res,
		adapter:      e.Adapter,
		handlers:     hf,
		index:        -1,
		pattern:      pattern,
		allowMethods: allowMethods,
	}
	c.Next()
}

// New create engine
func New[T any]() *Engine[T] {
	engine := &Engine[T]{
		Routes: NewRouter[T](),
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
