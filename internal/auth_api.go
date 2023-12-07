package internal

import (
	"github.com/goava/di"
	authapi "github.com/pridemon/outpost/pkg/auth_api"
	"github.com/pridemon/outpost/pkg/resty"
	"github.com/spf13/viper"
)

var AuthApiModule = di.Options(
	di.Provide(AuthApiConfigProvider),
	di.Provide(AuthApiProvider),
)

func AuthApiConfigProvider(v *viper.Viper) (*authapi.AuthApiConfig, error) {
	var config authapi.AuthApiConfig
	err := v.UnmarshalKey("auth_api", &config)
	return &config, err
}

func AuthApiProvider(restyConfig *resty.RestyConfig, config *authapi.AuthApiConfig) *authapi.AuthApi {
	return authapi.NewAuthApi(resty.RestyFactory(restyConfig), config)
}
