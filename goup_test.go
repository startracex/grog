package goup

import "testing"

func TestNew(t *testing.T) {
	e := New()
	e.Use(Recovery(), Logger())
	e.GET("/", func(request Request, response Response) {
		t.Log(request.URL)
		response.Status(200)
		response.String("Test New")
	})
	err := e.ListenAndServe("9527")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
