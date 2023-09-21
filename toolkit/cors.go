package toolkit

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrOriginNotMatch = errors.New("cors: origin not match")
var ErrOriginNotFound = errors.New("cors: origin not found")

type Cors struct {
	AllowOrigin      []string
	AllowMethod      []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int64
}

// CorsAllowAll Allow all
func CorsAllowAll() *Cors {
	return &Cors{
		AllowOrigin:      []string{"*"},
		AllowMethod:      []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"*"},
		MaxAge:           86400,
	}
}

// Match return if AllowOrigin match a header "Origin"
func (c *Cors) Match(origin string) bool {
	if len(c.AllowOrigin) == 0 {
		return false
	}
	if c.AllowOrigin[0] == "*" {
		return true
	}
	for _, allow := range c.AllowOrigin {
		if allow == origin {
			return true
		}
	}
	return false
}

// WriteHeaderOrigin Write to origin and other
func (c *Cors) WriteHeaderOrigin(header http.Header, origin string) {
	header.Set("Access-Control-Allow-Origin", origin)
	if len(c.AllowMethod) > 0 {
		header.Set("Access-Control-Allow-Methods", JoinString(c.AllowMethod))
	}
	if len(c.AllowHeaders) > 0 {
		header.Set("Access-Control-Allow-Headers", JoinString(c.AllowHeaders))
	}
	if len(c.ExposeHeaders) > 0 {
		header.Set("Access-Control-Expose-Headers", JoinString(c.ExposeHeaders))
	}
	if c.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.MaxAge > 0 {
		header.Set(
			"Access-Control-Max-Age", strconv.FormatInt(c.MaxAge, 10))
	}
}

// WriteHeader get "Origin" from getFrom, set Cors to writeTo
func (c *Cors) WriteHeader(getFrom, writeTo http.Header) error {
	origin := GetOrigin(getFrom)
	if origin == "" {
		return ErrOriginNotFound
	}
	if c.Match(origin) {
		c.WriteHeaderOrigin(writeTo, GetOrigin(getFrom))
		return nil
	}
	return ErrOriginNotMatch
}

// ToSecond convert time.Duration to int64 (Second)
func ToSecond(d time.Duration) int64 {
	return int64(d.Seconds())
}

// JoinString call strings.Join(s, ", ")
func JoinString(s []string) string {
	return strings.Join(s, ", ")
}

// GetOrigin get "Origin" from header
func GetOrigin(header http.Header) string {
	return header.Get("Origin")
}
