package goup

import (
    "testing"
)

func TestMethods(t *testing.T) {
    writeMethod := func(request Request, response Response) {
        response.String(request.Method)
    }
    e := New()
    e.METHOD(GET, "/get", writeMethod)
    e.NoRoute(writeMethod)
    e.ListenAndServe("9527")
}
