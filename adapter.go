package grog

import (
	"fmt"
	"net/http"
	"reflect"
)

func emptyAdapter(f HandlerFunc) HandlerFunc {
	return f
}

func httpAdapter(hf http.HandlerFunc) func(Context) {
	return func(ctx Context) {
		hf(ctx, ctx.Request())
	}
}

func defaultAdapter[T any](t T) func(Context) {
	i := reflect.ValueOf(t).Interface()
	if v, ok := i.(HandlerFunc); ok {
		return emptyAdapter(v)
	}
	if v, ok := i.(http.HandlerFunc); ok {
		return httpAdapter(v)
	}
	tName := reflect.TypeOf(t).String()
	panic(fmt.Sprintf("no supported adapter for %q", tName))
}
