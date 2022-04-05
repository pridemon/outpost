package proxy

import (
	"net/http"
	"strings"

	"github.com/goava/di"
	"github.com/sirupsen/logrus"
)

type ProxyConfig struct {
	Target string `json:"target" yaml:"target" mapstructure:"target"`
	Host   string `json:"host" yaml:"host" mapstructure:"host"`
}

type Proxy struct {
	di.Inject

	Log            *logrus.Logger
	WebProxy       http.Handler `di:"type=proxy_web"`
	WebsocketProxy http.Handler `di:"type=proxy_websocket"`
}

func (p *Proxy) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if p.isWebsocket(r) {
		p.Log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via websocket proxy")
		p.WebsocketProxy.ServeHTTP(w, r)
		return true
	}

	p.Log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via proxy")
	p.WebProxy.ServeHTTP(w, r)
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
