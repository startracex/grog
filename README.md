# grog

Grog is a minimalist web framework with domain/route grouping and pluggable handler adapters.

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
    hf(c.Writer(), c.Request())
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
