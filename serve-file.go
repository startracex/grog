package grog

import (
	"net/http"
	"os"
)

// ServeFile is similar to http.ServeFile, but it won't redirect
func ServeFile(c *Context, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	ServeContent(c, f)
}

// ServeContent serve file content
func ServeContent(c *Context, f *os.File) {
	s, err := f.Stat()
	if err != nil {
		return
	}
	http.ServeContent(c.Writer, c.Request, f.Name(), s.ModTime(), f)
}
