package reverse

import (
	"net/http"
	"net/http/httputil"
)

type ForwardInfo struct {
	Form    URL
	Target  URL
	BaseURL string
	Port    string
}

type ForwardFunc struct {
	BaseURL string
	Func    func(http.ResponseWriter, *http.Request)
}

type Forward struct {
	Form    *URL
	BaseURL string
	Target  *URL
}

// Add a new forward
func (e *Engine) Add(fw *Forward) {
	key := fw.Form.Host
	port := fw.Form.Port()
	if fw.Target.Scheme == "" {
		fw.Target.Scheme = "http"
	}

	// Insert information
	e.ListingPorts[port] = append(e.ListingPorts[port],
		ForwardInfo{
			Form:    *fw.Form,
			Target:  *fw.Target,
			BaseURL: fw.BaseURL,
			Port:    port,
		})

	e.ForwardFuncMaps[key] = append(e.ForwardFuncMaps[key], ForwardFunc{
		BaseURL: fw.BaseURL,
		Func: func(writer http.ResponseWriter, request *http.Request) {
			Reverse(fw.Target, writer, request)
		},
	})
}

// Reverse call NewSingleHostReverseProxy and ServeHTTP
func Reverse(url *URL, writer http.ResponseWriter, request *http.Request) {
	httputil.NewSingleHostReverseProxy(url).ServeHTTP(writer, request)
}
