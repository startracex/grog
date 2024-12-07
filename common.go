package goup

import (
	"github.com/startracex/goup/toolkit"
	"net/http"
)

type KV = map[string]any

// Upgrade to *toolkit.WS
func Upgrade(req Request, res Response) *toolkit.WS {
	return toolkit.Upgrade(res.Writer, req.OriginalRequest)
}

// Redirect call http.Redirect
func Redirect(request Request, response Response, url string, code int) {
	http.Redirect(response.Writer, request.OriginalRequest, url, code)
}

// WrapHandlerFunc wrap http.HandlerFunc to HandlerFunc
func WrapHandlerFunc(hf http.HandlerFunc) HandlerFunc {
	return func(request Request, response Response) {
		hf(response.Writer, request.OriginalRequest)
	}
}

// WrapHandler wrap http.Handler to HandlerFunc
func WrapHandler(h http.Handler) HandlerFunc {
	return func(request Request, response Response) {
		h.ServeHTTP(response.Writer, request.OriginalRequest)
	}
}

// ServeFile call http.ServeFile
func ServeFile(req Request, res Response, path string) {
	http.ServeFile(res.Writer, req.OriginalRequest, path)
}
