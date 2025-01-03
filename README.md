# goup

goup is a zero dependency web framework.

It provides APIs similar to `gin` and `express`.

## Import

```sh
go get -u github.com/startracex/goup
```

```go
import (
    "github.com/startracex/goup"
)
```

### Start a server

```go
engine := goup.New()
engine.GET("/", func (req goup.Request, res goup.Response) {
    res.String("Hello, world!")
})
engine.Run("9000")
```

```plain
Listen and serve at http://127.0.0.1:9000
```

### Middleware

```go
engine := goup.New()
engine.Use(/* ...middlewares */)
```

#### Defaults middlewares

```go
engine := goup.Default()
```

```go
engine := goup.New()
engine.Use(goup.DefaultMiddleware...)
```

### Router

#### Router group

```go
engine := goup.Default()
apiGroup := engine.Group("/api")
apiGroup.GET("/get-something", func (req Request, res Response) {

})
```

### Serve file

```go
engine := goup.New()
engine.Public("/favicon.ico", "./public/favicon.ico")
engine.Public("/public", "./public")
```

### WebSocket

```go
wsg := websocket.NewWSGroup()
engine.GET("/ws", func(request goup.Request, response goup.Response) {
    ws := goup.Upgrade(request, response)
    wsg.Add(ws)
    for {
        message := wsg.Message()
        if ws.Closed {
            break
        }
        wsg.Send(message, websocket.TEXT)
    }
})
```
