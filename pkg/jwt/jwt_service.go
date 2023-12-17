package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/goava/di"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

var (
	ErrBadAccessToken = errors.New("access token is bad")
)

type JwtConfig struct {
	SignKey         string        `json:"sign_key" yaml:"sign_key" mapstructure:"sign_key"`
	Issuer          string        `json:"iss" yaml:"iss" mapstructure:"iss"`
	Audience        string        `json:"aud" yaml:"aud" mapstructure:"aud"`
	RefreshTokenTTL time.Duration `json:"refresh_token_ttl" yaml:"refresh_token_ttl" mapstructure:"refresh_token_ttl"`
	WorkerDelay     time.Duration `json:"worker_delay" yaml:"worker_delay" mapstructure:"worker_delay"`
}

type JwtClaims struct {
	Subject uint   `json:"sub,omitempty"`
	Login   string `json:"login"`
	Email   string `json:"email"`
	jwt.StandardClaims
}

type JwtService struct {
	di.Inject

	Log    *logrus.Logger
	Config *JwtConfig
}

func (s *JwtService) CheckAccessToken(accessToken string) (*JwtClaims, error) {
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config.SignKey), nil
	})

	if err != nil || !token.Valid || claims.Issuer != s.Config.Issuer || claims.Audience != s.Config.Audience {
		err = fmt.Errorf("%w: %s", ErrBadAccessToken, err.Error())
		s.Log.Errorf("jwt.jwt_service: error checking access token: %v", err)
		return nil, err
	}

	return &claims, nil
}
