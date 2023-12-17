package internal

import (
	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/resty"
	"github.com/spf13/viper"
)

var RestyModule = di.Options(
	di.Provide(RestyConfigProvider),
)

func RestyConfigProvider(v *viper.Viper) (*resty.RestyConfig, error) {
	var config resty.RestyConfig
	err := v.UnmarshalKey("resty", &config)
	return &config, err
}
