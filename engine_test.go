package goup

import (
	"testing"
)

func TestRun(t *testing.T) {
	e := New()
	e.GET("/", func(request Request, response Response) {
		response.Status(200)
		response.String("Test New")
	})
	err := e.ListenAndServe("9527")
	if err != nil {
		t.Error(err)
	}
}

func TestRunTLS(t *testing.T) {
	e := New()
	e.GET("/", func(request Request, response Response) {
		response.Status(200)
		response.String("Test New TLS")
	})
	err := e.ListenAndServeTLS("9527", "server.crt", "server.key")
	if err != nil {
		t.Error(err)
	}
}
