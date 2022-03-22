package internal

import (
	"github.com/goava/di"
	"github.com/spf13/viper"
)

var HttpModule = di.Options(
	di.Provide(HttpConfigProvider),
)

type HttpConfig struct {
	Port int
}

func HttpConfigProvider(v *viper.Viper) (*HttpConfig, error) {
	var config HttpConfig
	err := v.UnmarshalKey("http", &config)
	return &config, err
}
