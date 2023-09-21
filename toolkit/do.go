package toolkit

import "net/http"

func Do(d *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(d)
}
