package grog

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Pattern  string
	Params   map[string]string
	Index    int
	Engine   *Engine
	Handlers []HandlerFunc
	Methods  []string
}

// Next call the next handler
func (c *Context) Next() {
	c.Index++
	for ; c.Index < len(c.Handlers); c.Index++ {
		c.Handlers[c.Index](c)
	}
}

// Abort handlers
func (c *Context) Abort() {
	c.Index = len(c.Handlers)
}

// Reset handlers
func (c *Context) Reset() {
	c.Index = -1
}
