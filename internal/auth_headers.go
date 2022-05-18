package internal

import (
	"github.com/goava/di"
	authheaders "github.com/pridemon/outpost/pkg/auth_headers"
	"github.com/spf13/viper"
)

var AuthHeadersModule = di.Options(
	di.Provide(AuthHeadersConfigProvider),
)

func AuthHeadersConfigProvider(v *viper.Viper) (*authheaders.AuthHeadersConfig, error) {
	var config authheaders.AuthHeadersConfig
	err := v.UnmarshalKey("auth_headers", &config)
	return &config, err
}
