package grog

// import (
// 	"net/http"

// 	"github.com/startracex/grog/websocket"
// )

// // // Upgrade wrap websocket.Upgrade
// // func Upgrade(c *Context) *websocket.WS {
// // 	return websocket.Upgrade(c.Writer, c.Request)
// // }

// // // Redirect wrap http.Redirect
// // func Redirect(c *Context, url string, code int) {
// // 	http.Redirect(c.Writer, c.Request, url, code)
// // }

// // // WrapHandlerFunc wrap http.HandlerFunc to HandlerFunc
// // func WrapHandlerFunc(hf http.HandlerFunc) HandlerFunc {
// // 	return func(c *Context) {
// // 		hf(c.Writer, c.Request)
// // 	}
// // }

// // // WrapHandler wrap http.Handler to HandlerFunc
// // func WrapHandler(h http.Handler) HandlerFunc {
// // 	return func(c *Context) {
// // 		h.ServeHTTP(c.Writer, c.Request)
// // 	}
// // }
