package grog

import (
	"net/http"
	"os"
	"time"
)

// ServeFile is similar to http.ServeFile, but it won't redirect
func ServeFile(c *Context, filePath string) {
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

// ServeContent serve file content
func ServeContent(c *Context, name string, modtime time.Time, f *os.File) {
	http.ServeContent(c.Writer, c.Request, name, modtime, f)
}
