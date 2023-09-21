package toolkit

import (
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := CorsAllowAll()
		match := c.Match(GetOrigin(writer.Header()))
		if match {
			c.WriteHeader(request.Header, writer.Header())
		}
	})
	http.ListenAndServe(":9527", nil)
}
