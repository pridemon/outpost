package authapi

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type AuthApiConfig struct {
	Host string `mapstructure:"host"`
}

type AuthApi struct {
	client *resty.Client
	config *AuthApiConfig
}

func NewAuthApi(restyClient *resty.Client, config *AuthApiConfig) *AuthApi {
	return &AuthApi{
		client: restyClient,
		config: config,
	}
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (p *AuthApi) Refresh(refreshToken string) (*RefreshResponse, error) {
	refreshUrl := p.config.Host + "/api/1.0/oauth/github/refresh"

	var respBody RefreshResponse

	resp, err := p.client.R().
		SetBody(map[string]interface{}{
			"refresh_token": refreshToken,
		}).
		SetResult(&respBody).
		Post(refreshUrl)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		err := fmt.Errorf("non 200 response from auth api: %d %s", resp.StatusCode(), resp.Body())
		return nil, err
	}

	return &respBody, nil
}
