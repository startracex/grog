# goup

goup is a simple web framework, Born from [gin](https://github.com/gin-gonic/gin/)
and [7days-golang](https://github.com/geektutu/7days-golang/), it provides usage similar to `(request, response) => { }`

## Run

### New

```go
engine := goup.New()
engine.GET("/", func (req Request, res Response) {
/* ... */
})
engine.Run(":9527")
```

### Use

```go
engine := goup.New()
engine.Use(func (req Request, res Response) {
/* ... */
})
engine.Use(Recovery(), Logger(), Cors())
```

### Group

```go
engine := goup.New()
api := engine.Group("/api")
api.GET("/xxx", func (req Request, res Response) {
/* ... */
})

```

### File

```go
engine := goup.New()
engine.File("/favicon.ico", "./public/favicon.ico")
engine.File("/public", "./public") //engine.Static("/public","./public")
```

### Quick Use

<details>
<summary>Request 
</summary>
get url

```go
req.URL()
```

get host

```go
req.Host()
```

get remote address

```go
req.Addr()
```

get path, params

```go
req.UseRouter()
```

get the key from params

```go
req.Param(key)
```

get URLSearchParams

```go
req.Query()
```

get key from URLSearchParams

```go
req.GetQuery()
```

get the key from form

```go
req.GetFormValue()
```

get the key file from form

```go
req.GetFormFile()
```

get all headers

```go
req.Header()
```

get the key from headers

```go
req.GetHeader(key)
```

get all cookies

```go
req.Cookies()
```

get the key from cookies

```go
req.GetCookie(key)
```

get body as string

```go
req.StringBody()
```

get body as bytes

```go
req.BytesBody()
```

#### context

set value

```go
req.Set(key, value)
req.SetValue(key, value)

```

get value

```go
req.Get(key)
req.GetValue(key)
```

</details>
<details>
<summary>Response
</summary>
writeheader(status)

```go
res.Status(200)
res.WriteHeader(200)
```

write

```go
res.Write([]byte("hello"))
res.Bytes([]byte("hello"))
res.String("hello")
```

header

```go
res.SetHeader("Content-Type", "text/html")
```

content type

```go
res.ContentType("text/html")
```

set cookie

```go
res.SetCookie(&cookie{})
```

html

```go
res.HTML("index.html", map[string]any{
"title": "goup",
})
```

json

```go
res.JSON(map[string]any{
"title": "goup",
})
```

status and error html page

```go
res.Error(404, "Not Found.")
```

</details>
