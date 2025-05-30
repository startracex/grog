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
engine := grog.New()
engine.GET("/", func (c *grog.Context) {
    c.Writer.Write([]byte("Hello World!"))
})
engine.Run("9000")
```

### With multiple domains

```go
engine := grog.New()
pagesDomain:=engine.Domain("pages.localhost")
pagesDomain.GET("/", func (c *grog.Context) {
    c.Writer.Write([]byte("Hello Pages!"))
})
apiDomain:=engine.Domain("api.localhost")
apiDomain.GET("/", func (c *grog.Context) {
    c.Writer.Write([]byte("Hello APIs!"))
})
engine.Run("9000")
```
