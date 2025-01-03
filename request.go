package goup

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type Request = *HttpRequest

type HandlerFunc func(Request, Response)

type HttpRequest struct {
	// original http request
	Reader   *http.Request
	Path     string
	Method   string
	Params   map[string]string
	Handlers []HandlerFunc
	index    int
	Engine   *Engine
}

func NewRequest(req *http.Request) HttpRequest {
	return HttpRequest{
		Reader: req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next call the next handler
func (r *HttpRequest) Next(w *HttpResponse) {
	r.index++
	for ; r.index < len(r.Handlers); r.index++ {
		r.Handlers[r.index](r, w)
	}
}

// Abort handlers
func (r *HttpRequest) Abort() {
	r.index = len(r.Handlers)
}

// Reset handlers
func (r *HttpRequest) Reset() {
	r.index = -1
}

// appendHandlers append handler functions
func (r *HttpRequest) appendHandlers(hs []HandlerFunc) {
	r.Handlers = append(r.Handlers, hs...)
}

/* Params */

// URL get url
func (r *HttpRequest) URL() *url.URL {
	return r.Reader.URL
}

// Host get host
func (r *HttpRequest) Host() string {
	return r.Reader.Host
}

// Addr return remote address
func (r *HttpRequest) Addr() string {
	return r.Reader.RemoteAddr
}

// Param get the key from params
func (r *HttpRequest) Param(key string) string {
	return r.Params[key]
}

// UseRouter get current path, params
func (r *HttpRequest) UseRouter() (string, map[string]string) {
	return r.Path, r.Params
}

// Query get URLSearchParams
func (r *HttpRequest) Query() url.Values {
	return r.Reader.URL.Query()
}

// GetQuery get key from URLSearchParams
func (r *HttpRequest) GetQuery(key string) string {
	return r.Reader.URL.Query().Get(key)
}

/* Form data */

// FormValue get the key from form
func (r *HttpRequest) FormValue(key string) string {
	return r.Reader.FormValue(key)
}

// PostFormValue get the key from form
func (r *HttpRequest) PostFormValue(key string) string {
	return r.Reader.PostFormValue(key)
}

// FormFile get the key file from form
func (r *HttpRequest) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return r.Reader.FormFile(key)
}

/* Data unmarshal */

// JSON unmarshal to v
func (r *HttpRequest) JSON(v any) error {
	return json.Unmarshal(r.BytesBody(), v)
}

func (r *HttpRequest) XML(v any) error {
	return xml.Unmarshal(r.BytesBody(), v)
}

/* Read */

// Body get *http.Request.Body
func (r *HttpRequest) Body() io.ReadCloser {
	return r.Reader.Body
}

// StringBody get body as buffer.String()
func (r *HttpRequest) StringBody() string {
	buf, _, err := r.copyBody()
	if err != nil {
		return ""
	}
	return buf.String()
}

// BytesBody get body as buffer.Bytes()
func (r *HttpRequest) BytesBody() []byte {
	buf, _, err := r.copyBody()
	if err != nil {
		return []byte{}
	}
	return buf.Bytes()
}

func (r *HttpRequest) copyBody() (*bytes.Buffer, int64, error) {
	buf := r.Engine.Pool.Get().(*bytes.Buffer)
	defer r.Engine.Pool.Put(buf)
	buf.Reset()
	l, err := io.Copy(buf, r.Reader.Body)
	return buf, l, err
}

/* Context */

// Context is alias of *http.Request.Context
func (r *HttpRequest) Context() context.Context {
	return r.Reader.Context()
}

// WithContext is alias of *http.Request.WithContext
func (r *HttpRequest) WithContext(ctx context.Context) {
	r.Reader = r.Reader.WithContext(ctx)
}

// SetValue Set custom parameters to the context
func (r *HttpRequest) SetValue(key any, value any) {
	r.WithContext(context.WithValue(r.Reader.Context(), key, value))
}

// Set is alias of SetValue
func (r *HttpRequest) Set(key string, value any) {
	r.SetValue(key, value)
}

// GetValue Get custom parameters to the context
func (r *HttpRequest) GetValue(key string) any {
	return r.Context().Value(key)
}

// Get is alias of GetValue
func (r *HttpRequest) Get(key string) any {
	return r.GetValue(key)
}

/* Cookies */

// Cookies get all cookies
func (r *HttpRequest) Cookies() []*http.Cookie {
	return r.Reader.Cookies()
}

// GetCookie get key from cookie
func (r *HttpRequest) GetCookie(key string) string {
	cookie, err := r.Reader.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

/* Headers */

// Header get header
func (r *HttpRequest) Header() http.Header {
	return r.Reader.Header
}

// GetHeader get the key from header
func (r *HttpRequest) GetHeader(key string) string {
	return r.Header().Get(key)
}

func (r *HttpRequest) ContentLength() int64 {
	return r.Reader.ContentLength
}

func (r *HttpRequest) ContentType() string {
	return r.GetHeader("Content-Type")
}

func (r *HttpRequest) Accept() string {
	return r.GetHeader("Accept")
}

func (r *HttpRequest) Authorization() string {
	return r.GetHeader("Authorization")
}

// BasicAuthorization set header "Authorization" with Basic scheme
func (r *HttpRequest) BasicAuthorization() string {
	a := r.Authorization()
	if strings.HasPrefix(a, "Basic ") {
		return a[6:]
	}
	return ""
}

// BearerAuthorization set header "Authorization" with Bearer scheme
func (r *HttpRequest) BearerAuthorization() string {
	a := r.Authorization()
	if strings.HasPrefix(a, "Bearer ") {
		return a[7:]
	}
	return ""
}

func (r *HttpRequest) UserAgent() string {
	return r.GetHeader("User-Agent")
}

func (r *HttpRequest) Origin() string {
	return r.GetHeader("Origin")
}
