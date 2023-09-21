package toolkit

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	var header = make(http.Header)
	NewCors().WriteHeader(&header, "*")
	fmt.Println(header)
}
