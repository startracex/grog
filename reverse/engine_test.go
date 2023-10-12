package reverse

import (
	"net/http"
	"testing"
)

const from = ":9527"
const to = ":9526"
const host = "localhost"

func TestRun(t *testing.T) {
	e := New()
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(to))
		})
		err := http.ListenAndServe(to, nil)
		if err != nil {
			t.Log(err)
		}
	}()
	e.AddOne(&Forward{
		Form:   &URL{Host: host + from},
		Target: &URL{Host: host + to},
	})
	e.Run()
}
