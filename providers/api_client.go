package providers

import "net/http"

type APIClient interface {
	Request(meth, endpoint string) (*http.Response, error)
}

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}
