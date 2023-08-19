package toolkit

import (
    "fmt"
    "net/http"
    "runtime"
    "testing"
    "time"
)

func TestWS_Message(t *testing.T) {
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ws := Upgrade(writer, request)
        for {
            msg, err := ws.Message()
            if err != nil {
                t.Log("Error:", err)
                return
            }
            fmt.Println(string(msg))
            ws.Send([]byte("Hello"), TEXT)
        }
    })
    http.ListenAndServe(":9527", nil)
}

func TestWS_Send(t *testing.T) {
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ws := Upgrade(writer, request)
        for range time.Tick(time.Second) {
            err := ws.Send([]byte("Hello"), TEXT)
            if err != nil {
                t.Log("Error:", err)
                return
            }
        }
    })
    http.ListenAndServe(":9527", nil)
}
func TestWSGroup_Send(t *testing.T) {
    wsg := WSGroup{}
    go func() {
        for range time.Tick(time.Second) {
            err := wsg.Send([]byte("Hello"), TEXT)
            if err != nil {
                t.Log("Error:", err)
            }
        }
    }()
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ws := Upgrade(writer, request)
        wsg.Add(ws)
    })
    http.ListenAndServe(":9527", nil)
}

func TestWSGroup_Message(t *testing.T) {
    wsg := WSGroup{}
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ws := Upgrade(writer, request)
        wsg.Add(ws)
        for {
            msg, err := ws.Message()
            if err != nil {
                t.Log("Error:", err)
                return
            }
            fmt.Println(string(msg), runtime.NumGoroutine())
        }
    })
    http.ListenAndServe(":9527", nil)
}
