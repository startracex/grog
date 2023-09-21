package goup

import "testing"

func TestMiddlewares(t *testing.T) {
	e := New()
	e.Use(Logger(2), Recovery(), Cors())
	e.GET("/", func(request Request, response Response) {
		panic("panic")
	})
	e.POST("/", func(request Request, response Response) {
		response.JSON(map[string]string{
			"url": "/",
		})
	})
	e.ListenAndServe("9527")
}
