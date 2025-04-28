package grog

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
)

// ServeFile is similar to http.ServeFile, but it won't redirect
func ServeFile(req Request, res Response, filePath string) {
	filePath = SecurePath(filePath)
	f, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			res.StatusText(404)
		} else if errors.Is(err, os.ErrPermission) {
			res.StatusText(403)
		} else {
			res.StatusText(500)
		}
		return
	}
	defer f.Close()
	ServeContent(req, res, f)
}

// ServeContent serve file content
func ServeContent(req Request, res Response, f *os.File) {
	s, err := f.Stat()
	if err != nil {
		res.StatusText(500)
		return
	}
	http.ServeContent(res.Writer, req.Reader, f.Name(), s.ModTime(), f)
}

// SecurePath ensures the security of the input path
func SecurePath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	return p
}
