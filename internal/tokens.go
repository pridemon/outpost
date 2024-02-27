package internal

import (
	"time"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/tokens"
	"github.com/spf13/viper"
)

var TokensModule = di.Options(
	di.Provide(TokensConfigProvider),
	di.Provide(tokens.NewTokensService),
)

func TokensConfigProvider(v *viper.Viper) (*tokens.TokensConfig, error) {
	// Set default values
	config := tokens.TokensConfig{
		CleanerDelay: 5 * time.Second,
	}

	err := v.UnmarshalKey("tokens_service", &config)
	return &config, err
}
