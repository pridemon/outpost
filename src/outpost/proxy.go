package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	Target *url.URL
	Host   string

	revProxy *httputil.ReverseProxy
}

func NewProxy() *Proxy {
	var proxy Proxy

	proxy.revProxy = &httputil.ReverseProxy{
		Director: proxy.Director,
	}

	return &proxy
}

func (p *Proxy) Director(req *http.Request) {
	req.URL.Scheme = p.Target.Scheme
	req.URL.Host = p.Target.Host

	if p.Host != "" {
		req.Host = p.Host
	}
}

func (p *Proxy) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via proxy")

	p.revProxy.ServeHTTP(w, r)
	return true
}
