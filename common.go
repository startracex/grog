package goup

import (
	"github.com/startracex/goup/toolkit"
	"net/http"
)

func Upgrade(req Request, res Response) *toolkit.WS {
	return toolkit.Upgrade(res.Writer, req.OriginalRequest)
}

func Do(d *http.Request) (*http.Response, error) {
	return toolkit.Do(d)
}
