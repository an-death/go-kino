package providers

import (
	"compress/gzip"
	"net/http"
	"net/url"
)

type APIClient interface {
	Request(meth, endpoint string) (*http.Response, error)
}

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

func NewProxyTransport(scheme, host string) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: scheme,
			Host:   host,
		}),
	}
}

func WrapGzipTransport(t http.RoundTripper) http.RoundTripper {
	return &GzipTransport{t}
}

type GzipTransport struct {
	http.RoundTripper
}

func (g *GzipTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Header.Get("Accept-Encoding") != "" {
		return g.RoundTripper.RoundTrip(r)
	}

	r.Header.Add("Accept-Encoding", "gzip")
	resp, err := g.RoundTripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
	}
	return resp, err
}
