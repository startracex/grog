package goup

import (
	"encoding/json"
	"fmt"
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

/* Status */

// With write header and call fn
// It helps you fill in the header and execute the callback fn
func (res HttpResponse) With(code int, fn func()) {
	res.WriteHeader(code)
	fn()
}

// StatusText send code http.StatusText(code)
func (res Response) StatusText(code int) {
	res.WriteHeader(code)
	res.String("%d %s", code, http.StatusText(code)+".")
}

// Status is alias of WriteHeader
func (res HttpResponse) Status(code int) HttpResponse {
	return res.WriteHeader(code)
}

// WriteHeader write status code
func (res HttpResponse) WriteHeader(code int) HttpResponse {
	res.Writer.WriteHeader(code)
	return res
}

/* Write */

// Write call Writer.Write
func (res HttpResponse) Write(data []byte) (int, error) {
	return res.Writer.Write(data)
}

// Byte is alias of Write
func (res HttpResponse) Byte(data []byte) (int, error) {
	return res.Write(data)
}

// String write string
func (res HttpResponse) String(format string, a ...any) (int, error) {
	if len(a) == 0 {
		return fmt.Fprint(res.Writer, format)
	}
	return fmt.Fprintf(res.Writer, format, a...)
}

/* Data encode */

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

/* Cookies */

// SetCookie set a cookie
func (res HttpResponse) SetCookie(cookie *http.Cookie) HttpResponse {
	http.SetCookie(res.Writer, cookie)
	return res
}

/* Headers */

// Header get header
func (res HttpResponse) Header() http.Header {
	return res.Writer.Header()
}

// SetHeader set a header
func (res HttpResponse) SetHeader(key, value string) HttpResponse {
	res.Header().Set(key, value)
	return res
}

// AddHeader add a header
func (res HttpResponse) AddHeader(key, value string) HttpResponse {
	res.Header().Add(key, value)
	return res
}

// DeleteHeader delete a header
func (res HttpResponse) DeleteHeader(key string) HttpResponse {
	res.Header().Del(key)
	return res
}

// Authorization set header "Authorization"
func (res HttpResponse) Authorization(scheme, parameters string) {
	res.SetHeader("Authorization", scheme+" "+parameters)
}

// BasicAuthorization set header "Authorization" with Basic scheme
func (res HttpResponse) BasicAuthorization(parameters string) {
	res.SetHeader("Authorization", "Basic: "+parameters)
}

// BearerAuthorization set header "Authorization" with Bearer scheme
func (res HttpResponse) BearerAuthorization(parameters string) {
	res.SetHeader("Authorization", "Bearer: "+parameters)
}

// ContentType set header "Content-Type"
func (res HttpResponse) ContentType(value string) {
	res.SetHeader("Content-Type", value)
}
