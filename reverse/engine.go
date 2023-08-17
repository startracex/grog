package reverse

import (
    "fmt"
    "net/http"
    "net/url"
)

type URL = url.URL

type Engine struct {
    ForwardFuncMaps map[string][]ForwardFunc
    ListingPorts    map[string][]ForwardInfo
}

func New() *Engine {
    return &Engine{
        ForwardFuncMaps: make(map[string][]ForwardFunc),
        ListingPorts:    make(map[string][]ForwardInfo),
    }
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    key := r.Host
    if forwardFuncs, ok := e.ForwardFuncMaps[key]; ok {
        maxlen := 0
        index := 0
        for _, f := range forwardFuncs {
            if f.BaseURL == "" || f.BaseURL == "/" {
                break
            }
            if r.URL.Path[:len(f.BaseURL)] == f.BaseURL {
                if len(f.BaseURL) > maxlen {
                    maxlen = len(f.BaseURL)
                    index++
                }
            }

        }
        forwardFuncs[index].Func(w, r)
    }
}

// Run Listen multiple address
func (e *Engine) Run(addr ...string) {

    for port, keys := range e.ListingPorts {
        fmt.Printf("Listening on %s:\n", port)
        for _, info := range keys {
            fmt.Printf("\t%s/%s => %s\n", info.Form.Host, info.BaseURL, info.Target.Host)
        }

        // append listening port to addr
        addr = append(addr, port)
    }

    for _, v := range addr {
        v := v
        go func() {
            err := e.ListenAndServe(v)
            if err != nil {
                fmt.Printf("Listening %s Error:%s/n", addr, err)
            }
        }()
    }
    select {}
}

func (e *Engine) ListenAndServe(addr string) error {
    if len(addr) > 0 && addr[0] != ':' {
        addr = ":" + addr
    }
    return http.ListenAndServe(addr, e)
}

func (e *Engine) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
    if len(addr) > 0 && addr[0] != ':' {
        addr = ":" + addr
    }
    return http.ListenAndServeTLS(addr, certFile, keyFile, e)
}
