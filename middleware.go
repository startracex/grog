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

type HandlerFunc = func(Context)

// Logger record the request path, method, time span.
func Logger() HandlerFunc {
	return func(c Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%s] %s <%v>", c.Method(), c.Path(), time.Since(t))
	}
}

var ErrRecovery = fmt.Errorf("%s", http.StatusText(500))

func Recovery() HandlerFunc {
	return func(c Context) {
		defer func() {
			err := recover()
			if err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.ResponseWriter().WriteHeader(500)
			}
		}()
		c.Next()
	}
}

// AutoOptions handle OPTIONS request, allow methods which have been registered.
func AutoOptions() HandlerFunc {
	return func(c Context) {
		config := &cors.Config{
			AllowOrigin:  []string{c.Request().Header.Get("Origin")},
			AllowMethods: c.AllowMethods(),
			AllowHeaders: []string{"*"},
			MaxAge:       86400,
		}
		config.WriteHeader(c.ResponseWriter().Header())
		if c.Request().Method == OPTIONS {
			c.ResponseWriter().WriteHeader(204)
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
