# grog

grog is a zero dependency web framework.

It provides APIs similar to `gin` and `express`.

## Import

```sh
go get -u github.com/startracex/grog
```

```go
import (
    "github.com/startracex/grog"
)
```

### Start a server

```go
engine := grog.New()
engine.GET("/", func (req grog.Request, res grog.Response) {
    res.String("Hello, world!")
})
engine.Run("9000")
```

```plain
Listen and serve at http://127.0.0.1:9000
```

### Middleware

```go
engine := grog.New()
engine.Use(/* ...middlewares */)
```

#### Defaults middlewares

```go
engine := grog.Default()
```

```go
engine := grog.New()
engine.Use(grog.DefaultMiddleware...)
```

### Router

#### Router group

```go
engine := grog.Default()
apiGroup := engine.Group("/api")
apiGroup.GET("/get-something", func (req Request, res Response) {

})
```

### Serve file

```go
engine := grog.New()
engine.Public("/favicon.ico", "./public/favicon.ico")
engine.Public("/public", "./public")
```

### WebSocket

```go
wsg := websocket.NewWSGroup()
engine.GET("/ws", func(request grog.Request, response grog.Response) {
    ws := grog.Upgrade(request, response)
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
