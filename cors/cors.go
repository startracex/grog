package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Config struct {
	Allow            []string
	AllowOrigin      []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int64
	RequestHeaders   []string
	RequestMethod    string
}

// MatchOrigin return if AllowOrigin match origin
func (c *Config) MatchOrigin(origin string) bool {
	for _, allow := range c.AllowOrigin {
		if allow == origin || allow == "*" {
			return true
		}
	}
	return false
}

func (c *Config) WriteHeader(header http.Header) {
	setHeaderValues(header, "Allow", c.AllowMethods)
	setHeaderValues(header, "Access-Control-Allow-Origin", c.AllowOrigin)
	setHeaderValues(header, "Access-Control-Allow-Methods", c.AllowMethods)
	setHeaderValues(header, "Access-Control-Allow-Headers", c.AllowHeaders)
	setHeaderValues(header, "Access-Control-Expose-Headers", c.ExposeHeaders)
	setHeaderValues(header, "Access-Control-Request-Headers", c.RequestHeaders)
	if c.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.MaxAge > 0 {
		header.Set("Access-Control-Max-Age", strconv.FormatInt(c.MaxAge, 10))
	}
	if c.RequestMethod != "" {
		header.Set("Access-Control-Request-Method", c.RequestMethod)
	}
}

// setHeaderValues set header values if values has length
func setHeaderValues(header http.Header, key string, values []string) {
	str := strings.Join(values, ", ")
	if str != "" {
		header.Set(key, str)
	}
}
