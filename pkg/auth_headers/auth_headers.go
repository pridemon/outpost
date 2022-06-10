package authheaders

import (
	"net/http"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type AuthHeadersConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	// key - property, value - headerName
	Headers map[string]string `json:"headers" yaml:"headers" mapstructure:"headers"`
}

type AuthHeadersService struct {
	di.Inject

	Log    *logrus.Logger
	Config *AuthHeadersConfig
}

func (s *AuthHeadersService) Process(r *http.Request, claims *jwt.JwtClaims) {
	if !s.Config.Enabled {
		return
	}

	values := map[string]string{
		"login": claims.Login,
		"email": claims.Email,
	}

	for property, headerName := range s.Config.Headers {
		r.Header.Set(headerName, values[property])
		s.Log.WithField("header", headerName).Debug("header is set")
	}
}
