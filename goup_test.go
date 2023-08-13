package goup

import (
    "testing"
)

func TestNew(t *testing.T) {
    e := New()
    e.GET("/", func(request Request, response Response) {
        response.Status(200)
        response.String("Test New")
    })
    e.ListenAndServe("9527")
}
