package goup

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Request = *HttpRequest

type HandlerFunc func(Request, Response)

type HttpRequest struct {
	// original http request
	OriginalRequest *http.Request
	Path            string
	Method          string
	Params          map[string]string
	Handlers        []HandlerFunc
	index           int
	Engine          *Engine
}

func NewRequest(req *http.Request) HttpRequest {
	return HttpRequest{
		OriginalRequest: req,
		Path:            req.URL.Path,
		Method:          req.Method,
		index:           -1,
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

/* Quick usage */

// URL get url
func (r *HttpRequest) URL() *url.URL {
	return r.OriginalRequest.URL
}

// Host get host
func (r *HttpRequest) Host() string {
	return r.OriginalRequest.Host
}

// Addr return remote address
func (r *HttpRequest) Addr() string {
	return r.OriginalRequest.RemoteAddr
}

// UseRouter get current path, params
func (r *HttpRequest) UseRouter() (string, map[string]string) {
	return r.Path, r.Params
}

// Param get the key from params
func (r *HttpRequest) Param(key string) string {
	return r.Params[key]
}

// Query get URLSearchParams
func (r *HttpRequest) Query() url.Values {
	return r.OriginalRequest.URL.Query()
}

// GetQuery get key from URLSearchParams
func (r *HttpRequest) GetQuery(key string) string {
	return r.OriginalRequest.URL.Query().Get(key)
}

// GetFormValue get the key from form
func (r *HttpRequest) GetFormValue(key string) string {
	return r.OriginalRequest.FormValue(key)
}

// GetFormFile get the key file from form
func (r *HttpRequest) GetFormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return r.OriginalRequest.FormFile(key)
}

// JSON unmarshal to v
func (r *HttpRequest) JSON(v any) error {
	return json.Unmarshal(r.BytesBody(), v)
}

// Header get header
func (r *HttpRequest) Header() http.Header {
	return r.OriginalRequest.Header
}

// GetHeader get the key from header
func (r *HttpRequest) GetHeader(key string) string {
	return r.Header().Get(key)
}

// Cookies get all cookies
func (r *HttpRequest) Cookies() []*http.Cookie {
	return r.OriginalRequest.Cookies()
}

// GetCookie get key from cookie
func (r *HttpRequest) GetCookie(key string) string {
	cookie, err := r.OriginalRequest.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// Body get *http.Request.Body
func (r *HttpRequest) Body() io.ReadCloser {
	return r.OriginalRequest.Body
}

// StringBody get body as buffer.String()
func (r *HttpRequest) StringBody() string {
	buf := r.Engine.Pool.Get().(*bytes.Buffer)
	defer r.Engine.Pool.Put(buf)
	buf.Reset()
	_, err := io.Copy(buf, r.OriginalRequest.Body)
	if err != nil {
		return ""
	}
	b := buf.String()
	return b
}

// BytesBody get body as buffer.Bytes()
func (r *HttpRequest) BytesBody() []byte {
	buf := r.Engine.Pool.Get().(*bytes.Buffer)
	r.Engine.Pool.Put(buf)
	buf.Reset()
	_, err := io.Copy(buf, r.OriginalRequest.Body)
	if err != nil {
		return []byte{}
	}
	b := buf.Bytes()
	return b
}

// Context is alias of *http.Request.Context
func (r *HttpRequest) Context() context.Context {
	return r.OriginalRequest.Context()
}

// WithContext is alias of *http.Request.WithContext
func (r *HttpRequest) WithContext(ctx context.Context) {
	r.OriginalRequest = r.OriginalRequest.WithContext(ctx)
}

// SetValue Set custom parameters to the context
func (r *HttpRequest) SetValue(key any, value any) {
	r.OriginalRequest = r.OriginalRequest.WithContext(
		context.WithValue(r.OriginalRequest.Context(), key, value),
	)
}

// Set is alias of SetValue
func (r *HttpRequest) Set(key string, value any) {
	r.SetValue(key, value)
}

// GetValue Get custom parameters to the context
func (r *HttpRequest) GetValue(key string) any {
	return r.OriginalRequest.Context().Value(key)
}

// Get is alias of GetValue
func (r *HttpRequest) Get(key string) any {
	return r.GetValue(key)
}
