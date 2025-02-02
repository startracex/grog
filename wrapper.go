package goup

import (
	"github.com/startracex/goup/websocket"
	"net/http"
)

// Upgrade wrap websocket.Upgrade
func Upgrade(req Request, res Response) *websocket.WS {
	return websocket.Upgrade(res.Writer, req.Reader)
}

// Redirect wrap http.Redirect
func Redirect(request Request, response Response, url string, code int) {
	http.Redirect(response.Writer, request.Reader, url, code)
}

// WrapHandlerFunc wrap http.HandlerFunc to HandlerFunc
func WrapHandlerFunc(hf http.HandlerFunc) HandlerFunc {
	return func(request Request, response Response) {
		hf(response.Writer, request.Reader)
	}
}

// WrapHandler wrap http.Handler to HandlerFunc
func WrapHandler(h http.Handler) HandlerFunc {
	return func(request Request, response Response) {
		h.ServeHTTP(response.Writer, request.Reader)
	}
}
