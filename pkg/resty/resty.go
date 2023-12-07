package resty

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type RestyConfig struct {
	RetryCount    int           `json:"retry_count" yaml:"retry_count" mapstructure:"retry_count"`
	RetryWaitTime time.Duration `json:"retry_wait_time" yaml:"retry_wait_time" mapstructure:"retry_wait_time"`
	Timeout       time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func RestyFactory(config *RestyConfig) *resty.Client {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")

	client.SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWaitTime)

	return client
}
