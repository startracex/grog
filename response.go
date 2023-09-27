package goup

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Response = *HttpResponse

type HttpResponse struct {
	Writer http.ResponseWriter
	Engine *Engine
}

func NewResponse(res http.ResponseWriter) HttpResponse {
	return HttpResponse{
		Writer: res,
	}
}

/* Quick usage */

// Status write status code
func (res HttpResponse) Status(code int) HttpResponse {
	return res.WriteHeader(code)
}

// WriteHeader write status code
func (res HttpResponse) WriteHeader(code int) HttpResponse {
	res.Writer.WriteHeader(code)
	return res
}

// Write call Writer.Write
func (res HttpResponse) Write(data []byte) (int, error) {
	return res.Writer.Write(data)
}

// Byte is alias of Write
func (res HttpResponse) Byte(data []byte) (int, error) {
	return res.Write(data)
}

// String write string
func (res HttpResponse) String(data string) (int, error) {
	return res.Write([]byte(data))
}

// Header get header
func (res HttpResponse) Header() http.Header {
	return res.Writer.Header()
}

// SetHeader set a header
func (res HttpResponse) SetHeader(key, value string) HttpResponse {
	res.Writer.Header().Set(key, value)
	return res
}

// ContentType set header "Content-Type"
func (res HttpResponse) ContentType(value string) {
	res.Writer.Header().Set("Content-Type", value)
}

// SetCookie set a cookie
func (res HttpResponse) SetCookie(cookie *http.Cookie) HttpResponse {
	http.SetCookie(res.Writer, cookie)
	return res
}

// JSON send JSON encoded data
func (res HttpResponse) JSON(data any) error {
	res.ContentType("application/json")
	encoder := json.NewEncoder(res.Writer)
	return encoder.Encode(data)
}

// HTML send HTML template
func (res HttpResponse) HTML(name string, data any) error {
	res.ContentType("text/html")
	return res.Engine.Template.ExecuteTemplate(res.Writer, name, data)
}

var ErrorHTML = `<title>{{.code}} {{.message}}</title><div style=height:100vh;text-align:center;display:flex;flex-direction:column;align-items:center;justify-content:center><div style=line-height:48px;height:48px><style>@media (prefers-color-scheme:light){body{color:#000;background:#fff;margin:0}}@media (prefers-color-scheme:dark){body{color:#fff;background:#000;margin:0}}</style><h1 style="display:inline-block;margin:0 20px 0 0;padding-right:22px;font-size:2em;font-weight:500;vertical-align:top;border-right:1px solid #808080">{{ .code }}</h1><h2 style="display:inline-block;margin:10px 0 10px 0;font-size:1.5em;font-weight:400;line-height:28px;vertical-align:top">{{.message}}</h2></div></div>`

// Error set status and send HTML with code, message
func (res Response) Error(code int, message string) {
	res.Status(code)
	res.ContentType("text/html")
	t, _ := template.New("goup").Parse(ErrorHTML)
	_ = t.Execute(res.Writer, map[string]any{
		"code":    code,
		"message": message,
	})
}

// ErrorStatusText call Error(code, http.StatusText(code)+".")
func (res Response) ErrorStatusText(code int) {
	res.Error(code, http.StatusText(code)+".")
}
