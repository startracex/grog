package goup

import (
	"encoding/json"
	"encoding/xml"
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
func (r Response) With(code int, fn func()) {
	r.WriteHeader(code)
	fn()
}

// StatusText send code http.StatusText(code)
func (r Response) StatusText(code int) {
	r.WriteHeader(code)
	r.String("%d %s", code, http.StatusText(code)+".")
}

// Status is alias of WriteHeader
func (r Response) Status(code int) Response {
	r.WriteHeader(code)
	return r
}

// WriteHeader write status code
func (r Response) WriteHeader(code int) {
	r.Writer.WriteHeader(code)
}

/* Write */

// Write call Writer.Write
func (r Response) Write(data []byte) (int, error) {
	return r.Writer.Write(data)
}

// Byte is alias of Write
func (r Response) Byte(data []byte) (int, error) {
	return r.Write(data)
}

// String write string
func (r Response) String(format string, a ...any) (int, error) {
	if len(a) == 0 {
		return fmt.Fprint(r.Writer, format)
	}
	return fmt.Fprintf(r.Writer, format, a...)
}

/* Data encode */

// JSON send JSON encoded data
func (r Response) JSON(data any) error {
	r.ContentType("application/json")
	encoder := json.NewEncoder(r.Writer)
	return encoder.Encode(data)
}

// XML send XML encoded data
func (r Response) XML(data any) error {
	r.ContentType("application/xml")
	encoder := xml.NewEncoder(r.Writer)
	return encoder.Encode(data)
}

// HTML send HTML template
func (r Response) HTML(name string, data any) error {
	r.ContentType("text/html")
	return r.Engine.Template.ExecuteTemplate(r.Writer, name, data)
}

/* Cookies */

// SetCookie set a cookie
func (r Response) SetCookie(cookie *http.Cookie) Response {
	http.SetCookie(r.Writer, cookie)
	return r
}

/* Headers */

// Header get header
func (r Response) Header() http.Header {
	return r.Writer.Header()
}

// SetHeader set a header
func (r Response) SetHeader(key, value string) Response {
	r.Header().Set(key, value)
	return r
}

// AddHeader add a header
func (r Response) AddHeader(key, value string) Response {
	r.Header().Add(key, value)
	return r
}

// DeleteHeader delete a header
func (r Response) DeleteHeader(key string) Response {
	r.Header().Del(key)
	return r
}

// Authorization set header "Authorization"
func (r Response) Authorization(scheme, parameters string) {
	r.SetHeader("Authorization", scheme+" "+parameters)
}

// BasicAuthorization set header "Authorization" with Basic scheme
func (r Response) BasicAuthorization(parameters string) {
	r.Authorization("Basic", parameters)
}

// BearerAuthorization set header "Authorization" with Bearer scheme
func (r Response) BearerAuthorization(parameters string) {
	r.Authorization("Bearer", parameters)
}

// ContentType set header "Content-Type"
func (r Response) ContentType(value string) {
	r.SetHeader("Content-Type", value)
}
