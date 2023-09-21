package toolkit

import (
	"net/http"
	"runtime"
)

func Attack(ducks string, n int) {
	req, err := http.NewRequest("GET", ducks, nil)
	if err != nil {
		panic(err)
	}
	AttackWithRequest(req, n)
}
func AttackWithRequest(ducks *http.Request, n int) {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	done := make(chan struct{})
	for i := 0; i < n; i++ {
		go func() {
			for {
				Do(ducks)
			}
		}()
	}
	<-done
}

func Do(d *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(d)
}
