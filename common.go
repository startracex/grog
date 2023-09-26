package goup

import (
	"github.com/startracex/goup/toolkit"
	"net/http"
)

// Upgrade to *toolkit.WS
func Upgrade(req Request, res Response) *toolkit.WS {
	return toolkit.Upgrade(res.Writer, req.OriginalRequest)
}

func Do(d *http.Request) (*http.Response, error) {
	return toolkit.Do(d)
}

// Redirect call http.Redirect
func Redirect(request Request, response Response, url string, code int) {
	http.Redirect(response.Writer, request.OriginalRequest, url, code)
}

// ServeFile call http.ServeFile
func ServeFile(req Request, res Response, path string) {
	http.ServeFile(res.Writer, req.OriginalRequest, path)
}
