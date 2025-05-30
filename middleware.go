package grog

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/startracex/grog/cors"
)

// var DefaultMiddleware = []HandlerFunc{Logger(), Recovery(), AutoOptions()}

// Logger record the request path, method
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%s] %s <%v>", c.Request.Method, c.Request.URL.Path, time.Since(t))
	}
}

var ErrRecovery = fmt.Errorf("%s", http.StatusText(500))

// Recovery error returns 500
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			err := recover()
			if err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Writer.WriteHeader(500)
			}
		}()
		c.Next()
	}
}

// Cors custom CORS config
func Cors(c *cors.Config) HandlerFunc {
	return func(ctx *Context) {
		c.WriteHeader(ctx.Writer.Header())
		ctx.Next()
	}
}

// AutoOptions handle OPTIONS request, allow methods which have been registered
func AutoOptions() HandlerFunc {
	return func(c *Context) {
		header := c.Writer.Header()
		header.Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))

		header.Set("Access-Control-Allow-Methods", strings.Join(c.Methods, ", "))
		header.Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == OPTIONS {
			c.Writer.WriteHeader(204)
			c.Abort()
			return
		}
		c.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(4, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
