package internal

import (
	"github.com/goava/di"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var LogrusModule = di.Options(
	di.Provide(LoggerConfigProvider),
	di.Provide(LogrusProvider),
)

type LoggerConfig struct {
	Format string `mapstructure:"format"`
}

func LoggerConfigProvider(v *viper.Viper) (*LoggerConfig, error) {
	var config LoggerConfig
	err := v.UnmarshalKey("log", &config)
	return &config, err
}

func LogrusProvider(config *LoggerConfig) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	if config.Format == "plain" {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	level, err := logrus.ParseLevel("debug")
	if err != nil {
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	return logger
}
