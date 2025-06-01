# grog

grog is a zero dependency web framework.

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
// Declare a engine
engine := grog.New[http.HandlerFunc]()

// Custom your adapter
engine.Adapter = func(hf http.HandlerFunc) func(grog.Context) {
  return func(c grog.Context) {
    hf(c.Writer(), c.Request())
  }
}

// Write handlers
engine.GET("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello, World!"))
})

// Start server
engine.Run("9000")
```

### With multiple domains

```go
pagesDomain := engine.Domain("pages.localhost")
pagesDomain.GET("/", func (c grog.Context) {
    c.Writer.Write([]byte("Hello Pages!"))
})

apiDomain := engine.Domain("api.localhost")
apiDomain.GET("/", func (c grog.Context) {
    c.Writer.Write([]byte("Hello APIs!"))
})

engine.Run("9000")
```
