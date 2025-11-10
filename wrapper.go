package grog

import (
	"io"
	"net/http"
	"os"
	"time"
)

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

// ServeFile is similar to http.ServeFile, but it won't redirect.
func ServeFile(c Context, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return
	}
	ServeContent(c, fi.Name(), fi.ModTime(), f)
}

// ServeContent serve file content.
func ServeContent(c Context, name string, modtime time.Time, f io.ReadSeeker) {
	http.ServeContent(c.Writer(), c.Request(), name, modtime, f)
}
