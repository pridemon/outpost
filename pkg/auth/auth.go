package auth

import (
	"html/template"
	"net/http"
	"regexp"

	"github.com/goava/di"
	authheaders "github.com/pridemon/outpost/pkg/auth_headers"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type AuthConfig struct {
	Title      string `json:"title" yaml:"title" mapstructure:"title"`
	Icon       string `json:"icon" yaml:"icon" mapstructure:"icon"`
	CookieName string `json:"cookie_name" yaml:"cookie_name" mapstructure:"cookie_name"`
	OAuthUrl   string `json:"oauth_url" yaml:"oauth_url" mapstructure:"oauth_url"`
}

type Auth struct {
	di.Inject

	Log                *logrus.Logger
	Config             *AuthConfig
	JwtService         *jwt.JwtService
	AuthHeadersService *authheaders.AuthHeadersService

	fserver http.Handler
	reURL   *regexp.Regexp
	page    *template.Template
}

func NewAuth(fserver http.Handler, reURL *regexp.Regexp, page *template.Template) *Auth {
	return &Auth{
		fserver: fserver,
		reURL:   reURL,
		page:    page,
	}
}

func (a *Auth) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if isValid, claims := a.checkCookie(r); isValid {
		a.AuthHeadersService.Process(r, claims)
		return false
	}

	a.Log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via fserver")

	if !a.reURL.MatchString(r.URL.Path) {
		a.serveAuthPage(w, r)
	} else {
		r.URL.Path = a.reURL.FindStringSubmatch(r.URL.Path)[1]
		a.fserver.ServeHTTP(w, r)
	}

	return true
}

func (a *Auth) checkCookie(r *http.Request) (bool, *jwt.JwtClaims) {
	cookie, err := r.Cookie(a.Config.CookieName)
	if err != nil {
		return false, nil
	}

	return a.JwtService.CheckAccessToken(cookie.Value)
}

func (a *Auth) serveAuthPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusUnauthorized)  // NOTE: 401 code helps js web-applications to determine expired tokens

	if err := a.page.Execute(w, a.Config); err != nil {
		a.Log.Error(err)
	}
}
