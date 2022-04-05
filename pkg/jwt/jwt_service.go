package jwt

import (
	"errors"

	"github.com/goava/di"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

var (
	ErrBadAccessToken = errors.New("access token is bad")
)

type JwtConfig struct {
	SignKey  string `json:"sign_key" yaml:"sign_key" mapstructure:"sign_key"`
	Issuer   string `json:"iss" yaml:"iss" mapstructure:"iss"`
	Audience string `json:"aud" yaml:"aud" mapstructure:"aud"`
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

func (s *JwtService) CheckAccessToken(accessToken string) bool {
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config.SignKey), nil
	})

	if err != nil || !token.Valid ||
		claims.Issuer != s.Config.Issuer || claims.Audience != s.Config.Audience {
		s.Log.WithField("error", ErrBadAccessToken)
		return false
	}

	return true
}
