package toolkit

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := CorsAllowAll()
		h := c.Match(GetOrigin(writer.Header()))
		c.WriteHeader(request.Header, writer.Header())
		fmt.Println(h)
	})
	http.ListenAndServe(":8080", nil)
}
