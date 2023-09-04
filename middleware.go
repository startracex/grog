package goup

import (
    "fmt"
    "log"
    "runtime"
    "strings"
    "time"
)

// Logger record the request path, method, custom time
func Logger(flag ...int) HandlerFunc {
    return func(req *HttpRequest, res *HttpResponse) {
        t := time.Now()
        req.Next(res)
        d := 0
        for _, v := range flag {
            d = d | v
        }
        log.SetFlags(d)
        log.Printf("[%s] %s <%v>", req.Method, req.Path, time.Since(t))
    }
}

// Recovery error returns 500
func Recovery() HandlerFunc {
    return func(req *HttpRequest, res *HttpResponse) {
        defer func() {
            if err := recover(); err != nil {
                message := fmt.Sprintf("%s", err)
                log.Printf("%s\n\n", trace(message))
                res.Error(500, "INTERNAL SERVER ERROR.")
            }
        }()
        req.Next(res)
    }
}

// Cors (Allow-Origin?, Allow-Methods?, Allow-Headers?)
func Cors(s ...string) HandlerFunc {
    allows := []string{"*", "*", "*"}
    for i, v := range s {
        allows[i] = v
    }
    allowOrigin := allows[0]
    allowMethods := allows[1]
    allowHeaders := allows[2]
    return func(req Request, res Response) {
        res.SetHeader("Access-Control-Allow-Origin", allowOrigin)
        res.SetHeader("Access-Control-Allow-Methods", allowMethods)
        res.SetHeader("Access-Control-Allow-Headers", allowHeaders)
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
