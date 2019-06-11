package clients

import "net/http"

type APIClient interface {
	Request(meth, baseUrl, endpoint string) (*http.Response, error)
}

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}
