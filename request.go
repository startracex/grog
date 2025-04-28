package grog

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

type Request = *InnerRequest

type HandlerFunc func(Request, Response)

type InnerRequest struct {
	// original http request
	Reader   *http.Request
	Path     string
	Pattern  string
	Method   string
	Params   map[string]string
	Body     io.ReadCloser
	Header   http.Header
	Query    url.Values
	Cookies  []*http.Cookie
	Handlers []HandlerFunc
	index    int
	Engine   *Engine
}

func NewRequest(req *http.Request) InnerRequest {
	return InnerRequest{
		Reader:  req,
		Path:    req.URL.Path,
		Method:  req.Method,
		Body:    req.Body,
		Query:   req.URL.Query(),
		Header:  req.Header,
		Cookies: req.Cookies(),
		index:   -1,
	}
}

// Next call the next handler
func (r Request) Next(w *InnerResponse) {
	r.index++
	for ; r.index < len(r.Handlers); r.index++ {
		r.Handlers[r.index](r, w)
	}
}

// Abort handlers
func (r Request) Abort() {
	r.index = len(r.Handlers)
}

// Reset handlers
func (r Request) Reset() {
	r.index = -1
}

// appendHandlers append handler functions
func (r Request) appendHandlers(hs []HandlerFunc) {
	r.Handlers = append(r.Handlers, hs...)
}

/* Params */

// URL get url
func (r Request) URL() *url.URL {
	return r.Reader.URL
}

// Host get host
func (r Request) Host() string {
	return r.Reader.Host
}

// Addr return remote address
func (r Request) Addr() string {
	return r.Reader.RemoteAddr
}

// Param get the key from params
func (r Request) Param(key string) string {
	return r.Params[key]
}

// UseRouter get current path, params
func (r Request) UseRouter() (string, map[string]string) {
	return r.Path, r.Params
}

// GetQuery get key from URLSearchParams
func (r Request) GetQuery(key string) string {
	return r.Query.Get(key)
}

// HasQuery check if key exists in URLSearchParams
func (r Request) HasQuery(key string) bool {
	return r.Query.Has(key)
}

/* Form data */

// FormValue get the key from form
func (r Request) FormValue(key string) string {
	return r.Reader.FormValue(key)
}

// PostFormValue get the key from form
func (r Request) PostFormValue(key string) string {
	return r.Reader.PostFormValue(key)
}

// FormFile get the key file from form
func (r Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return r.Reader.FormFile(key)
}

/* Data unmarshal */

// JSON unmarshal to v
func (r Request) JSON(v any) error {
	return json.Unmarshal(r.BytesBody(), v)
}

func (r Request) XML(v any) error {
	return xml.Unmarshal(r.BytesBody(), v)
}

/* Read */

// StringBody get body as buffer.String()
func (r Request) StringBody() string {
	buf, _, err := r.copyBody()
	if err != nil {
		return ""
	}
	return buf.String()
}

// BytesBody get body as buffer.Bytes()
func (r Request) BytesBody() []byte {
	buf, _, err := r.copyBody()
	if err != nil {
		return []byte{}
	}
	return buf.Bytes()
}

func (r Request) copyBody() (*bytes.Buffer, int64, error) {
	buf := r.Engine.Pool.Get().(*bytes.Buffer)
	defer r.Engine.Pool.Put(buf)
	buf.Reset()
	l, err := io.Copy(buf, r.Body)
	return buf, l, err
}

/* Context */

// WithContext is alias of *http.Request.WithContext
func (r Request) WithContext(ctx context.Context) {
	r.Reader = r.Reader.WithContext(ctx)
}

// SetValue Set custom parameters to the context
func (r Request) SetValue(key any, value any) {
	r.WithContext(context.WithValue(r.Reader.Context(), key, value))
}

// Set is alias of SetValue
func (r Request) Set(key string, value any) {
	r.SetValue(key, value)
}

// GetValue Get custom parameters to the context
func (r Request) GetValue(key string) any {
	return r.Reader.Context().Value(key)
}

// Get is alias of GetValue
func (r Request) Get(key string) any {
	return r.GetValue(key)
}

/* Cookies */

// GetCookie get key from cookie
func (r Request) GetCookie(key string) string {
	cookie, err := r.Reader.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

/* Headers */

// GetHeader get the key from header
func (r Request) GetHeader(key string) string {
	return r.Header.Get(key)
}

func (r Request) ContentLength() int64 {
	return r.Reader.ContentLength
}

func (r Request) ContentType() string {
	return r.GetHeader("Content-Type")
}

func (r Request) Accept() string {
	return r.GetHeader("Accept")
}

func (r Request) Authorization() string {
	return r.GetHeader("Authorization")
}

// BasicAuthorization set header "Authorization" with Basic scheme
func (r Request) BasicAuthorization() string {
	a := r.Authorization()
	if strings.HasPrefix(a, "Basic ") {
		return a[6:]
	}
	return ""
}

// BearerAuthorization set header "Authorization" with Bearer scheme
func (r Request) BearerAuthorization() string {
	a := r.Authorization()
	if strings.HasPrefix(a, "Bearer ") {
		return a[7:]
	}
	return ""
}

func (r Request) UserAgent() string {
	return r.GetHeader("User-Agent")
}

func (r Request) Origin() string {
	return r.GetHeader("Origin")
}
