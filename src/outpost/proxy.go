package main

import (
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	webProxy       http.Handler
	websocketProxy http.Handler
}

func NewProxy(target, host string) *Proxy {
	var proxy Proxy

	turl, err := url.Parse(target)
	if err != nil {
		log.WithError(err).Fatal("OUTPOST_TARGET is not valid url")
	}

	proxy.webProxy = NewWebProxy(turl, host)
	proxy.websocketProxy = NewWebsocketProxy(turl.Host)

	return &proxy
}

func (p *Proxy) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if p.isWebsocket(r) {
		log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via websocket proxy")
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
