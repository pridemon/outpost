package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/koding/websocketproxy"
	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	Target *url.URL
	Host   string

	webProxy       *httputil.ReverseProxy
	websocketProxy *websocketproxy.WebsocketProxy
}

func NewProxy() *Proxy {
	var proxy Proxy

	proxy.webProxy = &httputil.ReverseProxy{
		Director: proxy.director,
	}

	proxy.websocketProxy = &websocketproxy.WebsocketProxy{
		Backend: proxy.websocketBackend,
	}

	return &proxy
}

func (p *Proxy) director(req *http.Request) {
	req.URL.Scheme = p.Target.Scheme
	req.URL.Host = p.Target.Host

	if p.Host != "" {
		req.Host = p.Host
	}
}

func (p *Proxy) websocketBackend(req *http.Request) *url.URL {
	// shallow copy
	url := *(req.URL)
	url.Host = p.Target.Host
	url.Scheme = p.Target.Scheme
	return &url
}

func (p *Proxy) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if p.isWebsocket(r) {
		log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via websocket")
		p.websocketProxy.ServeHTTP(w, r)
		return true
	}

	log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via proxy")
	p.webProxy.ServeHTTP(w, r)
	return true
}

func (p *Proxy) isWebsocket(req *http.Request) bool {
	res := false
	if strings.ToLower(p.getFirstHeader(req, "Connection")) == "upgrade" {
		res = p.getFirstHeader(req, "Upgrade") == "websocket"
	}

	return res
}

func (p *Proxy) getFirstHeader(req *http.Request, name string) string {
	headers := req.Header[name]
	if len(headers) == 0 {
		return ""
	}
	return headers[0]
}
