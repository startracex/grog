package grog

import (
	"net/http"
)

type ErrorBuilder interface {
	error
	Build(http.ResponseWriter)
}
