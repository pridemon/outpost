package internal

import (
	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/spf13/viper"
)

var JwtModule = di.Options(
	di.Provide(JwtConfigProvider),
)

func JwtConfigProvider(v *viper.Viper) (*jwt.JwtConfig, error) {
	var config jwt.JwtConfig
	err := v.UnmarshalKey("jwt", &config)
	return &config, err
}
