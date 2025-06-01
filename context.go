package grog

import (
	"net/http"

	"github.com/startracex/grog/router"
)

type Context interface {
	Writer() http.ResponseWriter
	Request() *http.Request
	Next()
	Abort()
	Reset()
}

type HandleContext[T any] struct {
	request        *http.Request
	writer         http.ResponseWriter
	path           string
	pattern        string
	params         map[string]string
	index          int
	HandlerAdapter func(T) func(Context)
	handlers       []T
	allowMethods   []string
}

// Next call the next handler
func (c *HandleContext[T]) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		fn := c.HandlerAdapter(c.handlers[c.index])
		if fn != nil {
			fn(c)
		}
	}
}

// Abort handlers
func (c *HandleContext[T]) Abort() {
	c.index = len(c.handlers)
}

// Reset handlers
func (c *HandleContext[T]) Reset() {
	c.index = -1
}

func (c *HandleContext[T]) Request() *http.Request {
	return c.request
}

func (c *HandleContext[T]) Writer() http.ResponseWriter {
	return c.writer
}

func (c *HandleContext[T]) Params() map[string]string {
	if c.params == nil {
		c.params = router.ParseParams(c.pattern, c.Path())
	}
	return c.params
}

func (c *HandleContext[T]) Pattern() string {
	return c.pattern
}

func (c *HandleContext[T]) AllowMethods() []string {
	return c.allowMethods
}

func (c *HandleContext[T]) Path() string {
	return c.path
}
