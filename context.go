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
}

type HandleContext[T any] struct {
	request        *http.Request
	writer         http.ResponseWriter
	Pattern        string
	Params         map[string]string
	Index          int
	HandlerAdapter func(T) func(Context)
	Handlers       []T
	Methods        []string
}

// Next call the next handler
func (c *HandleContext[T]) Next() {
	c.Index++
	for ; c.Index < len(c.Handlers); c.Index++ {
		fn := c.HandlerAdapter(c.Handlers[c.Index])
		if fn != nil {
			fn(c)
		}
	}
}

// Abort handlers
func (c *HandleContext[T]) Abort() {
	c.Index = len(c.Handlers)
}

// Reset handlers
func (c *HandleContext[T]) Reset() {
	c.Index = -1
}

func (c *HandleContext[T]) Request() *http.Request {
	return c.request
}

func (c *HandleContext[T]) Writer() http.ResponseWriter {
	return c.writer
}
