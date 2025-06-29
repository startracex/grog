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
}

func AllowAll() *Config {
	return &Config{
		AllowOrigin:   []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
		MaxAge:        86400,
	}
}

// MatchOrigin return if AllowOrigin match origin
func (c *Config) MatchOrigin(origin string) bool {
	if c.AllowOrigin == nil {
		return false
	}
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
	if c.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.MaxAge > 0 {
		header.Set(
			"Access-Control-Max-Age", strconv.FormatInt(c.MaxAge, 10))
	}
}

// setHeaderValues set header values if values has length
func setHeaderValues(header http.Header, key string, values []string) {
	if values == nil {
		return
	}
	str := strings.Join(values, ", ")
	if str != "" {
		header.Set(key, str)
	}
}
