package grog

import (
	"net/http"
	"os"
	"time"
)

// ServeFile is similar to http.ServeFile, but it won't redirect
func ServeFile(c *HandleContext[any], filePath string) {
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
func ServeContent(c *HandleContext[any], name string, modtime time.Time, f *os.File) {
	http.ServeContent(c.writer, c.request, name, modtime, f)
}
