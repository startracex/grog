package toolkit

import (
	"net/http"
	"testing"
)

func TestAttackWithRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:9527/", nil)
	AttackWithRequest(request, 1)
}
