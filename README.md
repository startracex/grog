# grog

Grog is a minimalist web framework with domain-based routing and pluggable handler adapters.

The Grog framework well-suited for developers who want a lightweight, extensible foundation without the overhead of heavier frameworks, or for beginners looking to learn how web frameworks work.

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
// Declare a engine with handlers type.
engine := grog.New[http.HandlerFunc]()

// Custom your adapter, for http.HandlerFunc and grog.HandlerFunc,
// this step can be omitted, the default converter will be used.
engine.Adapter = func(hf http.HandlerFunc) func(grog.Context) {
  return func(c grog.Context) {
    hf(c.ResponseWriter(), c.Request())
  }
}

// Write handlers, the first parameter is pattern,
// and the remaining parameters are handlers,
// the type of handler is generic when the engine is created.
engine.GET("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello, World!"))
})

// Start server.
engine.Run("9000")
```

### With multiple domains

```go
package main

import (
  "github.com/startracex/grog"
  "github.com/startracex/grog/websocket"
)

func main() {
	engine := grog.New[grog.HandlerFunc]()

  pagesDomain := engine.Domain("pages.localhost")
  pagesDomain.GET("/", func (c grog.Context) {
    c.ResponseWriter().Write([]byte("Hello Pages!"))
  })

  apiDomain := engine.Domain("api.localhost")
  apiDomain.GET("/", func (c grog.Context) {
    c.ResponseWriter().Write([]byte("Hello APIs!"))
  })

  chatDomain := engine.Domain("chat.localhost")
  chatDomain.GET("/", func (c grog.Context) {
    ws := websocket.New()
    ws.Upgrade(c.ResponseWriter(), c.Request())
    for {
      data, datatype, err := ws.Message()
      if ws.Closed || err != nil {
        break
      }
      ws.Send(data, datatype)
    }
  })

  engine.Run("9000")
}
```
