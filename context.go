package grog

import (
	"bufio"
	"net"
	"net/http"

	"github.com/startracex/grog/router"
)

type Context interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Next()
	Abort()
	Reset()
	Pattern() string
	Path() string
	Params() map[string]string
	Method() string
	AllowMethods() []string
}

type handleContext[T any] struct {
	request  *http.Request
	response http.ResponseWriter
	pattern  string
	index    int
	adapter  func(T) func(Context)
	handlers []T
	node     *router.Router[map[string][]T]
}

// Next call the next handler
func (c *handleContext[T]) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		fn := c.adapter(c.handlers[c.index])
		if fn != nil {
			fn(c)
		}
	}
}

// Abort handlers
func (c *handleContext[T]) Abort() {
	c.index = len(c.handlers)
}

// Reset handlers
func (c *handleContext[T]) Reset() {
	c.index = -1
}

func (c *handleContext[T]) Request() *http.Request {
	return c.request
}

func (c *handleContext[T]) Pattern() string {
	return c.pattern
}

func (c *handleContext[T]) Path() string {
	return c.Request().URL.Path
}

func (c *handleContext[T]) Params() map[string]string {
	return router.ParseParams(c.Path(), c.Pattern())
}

func (c *handleContext[T]) Method() string {
	return c.Request().Method
}

func (c *handleContext[T]) AllowMethods() []string {
	allowMethods := make([]string, len(c.node.Value))
	for method := range c.node.Value {
		allowMethods = append(allowMethods, method)
	}
	return allowMethods
}

func (c *handleContext[T]) ResponseWriter() http.ResponseWriter {
	return c.response
}

func (c *handleContext[T]) Header() http.Header {
	return c.response.Header()
}

func (c *handleContext[T]) Write(b []byte) (int, error) {
	return c.response.Write(b)
}

func (c *handleContext[T]) WriteHeader(statusCode int) {
	c.response.WriteHeader(statusCode)
}

func (c *handleContext[T]) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return c.response.(http.Hijacker).Hijack()
}

func (c *handleContext[T]) Flush() {
	c.response.(http.Flusher).Flush()
}
