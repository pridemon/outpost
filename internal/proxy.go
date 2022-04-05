package internal

import (
	"net/http"
	"net/url"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/proxy"
	"github.com/spf13/viper"
)

var ProxyModule = di.Options(
	di.Provide(ProxyConfigProvider),
	di.Provide(ProxyWebProvider, di.Tags{"type": "proxy_web"}),
	di.Provide(ProxyWebSocketProvider, di.Tags{"type": "proxy_websocket"}),
)

func ProxyConfigProvider(v *viper.Viper) (*proxy.ProxyConfig, error) {
	var config proxy.ProxyConfig
	err := v.UnmarshalKey("proxy", &config)
	return &config, err
}

func ProxyWebProvider(config *proxy.ProxyConfig) (http.Handler, error) {
	turl, err := url.Parse(config.Target)
	if err != nil {
		return nil, err
	}
	return proxy.NewWebProxy(turl, config.Host), nil
}

func ProxyWebSocketProvider(config *proxy.ProxyConfig) (http.Handler, error) {
	turl, err := url.Parse(config.Target)
	if err != nil {
		return nil, err
	}
	return proxy.NewWebsocketProxy(turl.Host), nil
}
