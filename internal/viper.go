package internal

import (
	"strings"

	"github.com/goava/di"
	"github.com/spf13/viper"
)

var ViperModule = di.Options(
	di.Provide(ViperProvider),
)

const (
	ViperEnvPrefix  = "OUTPOST"
	ViperConfigName = "config"
	ViperConfigPath = "/etc/outpost"
)

func ViperProvider() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(ViperEnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigName(ViperConfigName)
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath(ViperConfigPath)
	return v, v.ReadInConfig()
}
