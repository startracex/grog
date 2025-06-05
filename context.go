package grog

import (
	"net/http"
)

type Context interface {
	Writer() http.ResponseWriter
	Request() *http.Request
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
	request      *http.Request
	writer       http.ResponseWriter
	pattern      string
	params       map[string]string
	index        int
	adapter      func(T) func(Context)
	handlers     []T
	allowMethods []string
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

func (c *handleContext[T]) Writer() http.ResponseWriter {
	return c.writer
}

func (c *handleContext[T]) Pattern() string {
	return c.pattern
}

func (c *handleContext[T]) Path() string {
	return c.Request().URL.Path
}

func (c *handleContext[T]) Params() map[string]string {
	return c.params
}

func (c *handleContext[T]) Method() string {
	return c.Request().Method
}

func (c *handleContext[T]) AllowMethods() []string {
	return c.allowMethods
}
