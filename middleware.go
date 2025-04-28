package grog

import (
	"fmt"
	"github.com/startracex/grog/cors"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var DefaultMiddleware = []HandlerFunc{Logger(), Recovery(), AutoOptions()}

// Logger record the request path, method
func Logger() HandlerFunc {
	return func(req *InnerRequest, res *InnerResponse) {
		t := time.Now()
		req.Next(res)
		log.Printf("[%s] %s <%v>", req.Method, req.Path, time.Since(t))
	}
}

var ErrRecovery = fmt.Errorf(http.StatusText(500))

// Recovery error returns 500
func Recovery() HandlerFunc {
	return func(req *InnerRequest, res *InnerResponse) {
		defer func() {
			err := recover()
			if err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				res.Error(ErrRecovery)
			}
		}()
		req.Next(res)
	}
}

// Cors custom CORS config
func Cors(c *cors.Config) HandlerFunc {
	return func(req Request, res Response) {
		c.WriteHeader(res.Header())
		req.Next(res)
	}
}

// AutoOptions handle OPTIONS request, allow methods which have been registered
func AutoOptions() HandlerFunc {
	return func(req Request, res Response) {
		res.SetHeader("Access-Control-Allow-Origin", req.Origin())
		if req.Method == OPTIONS {
			methods := strings.Join(req.Engine.router.handlers.allMethods(req.Pattern), ", ")
			if methods == "" {
				res.WriteHeader(404)
			} else {
				res.SetHeader("Access-Control-Allow-Methods", methods)
				res.SetHeader("Access-Control-Allow-Headers", "*")
				res.SetHeader("Access-Control-Allow-Credentials", "true")
				res.WriteHeader(204)
			}
			req.Abort()
			return
		}
		req.Next(res)
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
