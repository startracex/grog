package toolkit

import (
    "fmt"
    "net/http"
    "testing"
)

func TestUpdateToWebsocket(t *testing.T) {
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ws := Upgrade(writer, request)
        i := 1
        for {
            data, _ := ws.Message()
            fmt.Println(i, string(data))
            _ = ws.Send(data, TEXT)
            i++
        }

    })
    http.ListenAndServe(":9527", nil)
}
