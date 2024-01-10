package internal

import (
	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/tokens"
)

var TokensModule = di.Options(
	di.Provide(tokens.NewTokensService),
)
