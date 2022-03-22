package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewWebProxy(turl *url.URL, host string) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = turl.Scheme
			req.URL.Host = turl.Host

			if host != "" {
				req.Host = host
			}
		},
	}
}
