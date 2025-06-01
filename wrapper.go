package grog

import (
	"net/http"

	"github.com/startracex/grog/websocket"
)

func Upgrade(c Context) *websocket.WS {
	return websocket.Upgrade(c.Writer(), c.Request())
}

func Redirect(c Context, url string, code int) {
	http.Redirect(c.Writer(), c.Request(), url, code)
}

func WrapHandlerFunc(hf http.HandlerFunc) HandlerFunc {
	return func(c Context) {
		hf(c.Writer(), c.Request())
	}
}

func WrapHandler(h http.Handler) HandlerFunc {
	return func(c Context) {
		h.ServeHTTP(c.Writer(), c.Request())
	}
}
