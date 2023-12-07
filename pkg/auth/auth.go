package auth

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"

	"github.com/goava/di"
	authheaders "github.com/pridemon/outpost/pkg/auth_headers"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/pridemon/outpost/pkg/tokens"
	"github.com/sirupsen/logrus"
)

type AuthConfig struct {
	Title             string `json:"title" yaml:"title" mapstructure:"title"`
	Icon              string `json:"icon" yaml:"icon" mapstructure:"icon"`
	AccessCookieName  string `json:"access_cookie_name" yaml:"access_cookie_name" mapstructure:"access_cookie_name"`
	RefreshCookieName string `json:"refresh_cookie_name" yaml:"refresh_cookie_name" mapstructure:"refresh_cookie_name"`
	CookieDomain      string `json:"cookie_domain" yaml:"cookie_domain" mapstructure:"cookie_domain"`
	OAuthUrl          string `json:"oauth_url" yaml:"oauth_url" mapstructure:"oauth_url"`
}

type Auth struct {
	di.Inject

	Log                *logrus.Logger
	Config             *AuthConfig
	AuthHeadersService *authheaders.AuthHeadersService
	TokensService      *tokens.TokensService

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
	if a.reURL.MatchString(r.URL.Path) {
		return a.serveViaFserver(w, r)
	}

	var accessCookieValue, refreshCookieValue string

	accessCookie, err := r.Cookie(a.Config.AccessCookieName)
	if err != nil {
		a.Log.Errorf("auth: error getting cookie with access token: %v", err)
		return a.serveAuthPage(w, r)
	}
	accessCookieValue = accessCookie.Value

	refreshCookie, _ := r.Cookie(a.Config.RefreshCookieName)
	if refreshCookie != nil && refreshCookie.Value != "" {
		refreshCookieValue = refreshCookie.Value
		a.deleteCookie(w, refreshCookie)
	}

	claims, err := a.TokensService.ProcessTokens(accessCookieValue, refreshCookieValue)

	if errors.Is(err, jwt.ErrBadAccessToken) && !a.reURL.MatchString(r.URL.Path) {
		a.Log.WithField("error", err).Debug("auth: trying to refresh access token")

		accessToken, err := a.TokensService.RefreshToken(accessCookieValue)
		if err != nil {
			a.Log.Errorf("auth: error refreshing access token: %v", err)
			return a.serveAuthPage(w, r)
		}

		claims, err = a.TokensService.ProcessAccessToken(accessToken)
		if err != nil {
			a.Log.Errorf("auth: error processing access token: %v", err)
			return a.serveAuthPage(w, r)
		}

		a.updateCookie(w, accessCookie, accessToken)
	} else if err != nil {
		a.Log.Errorf("auth: error processing tokens: %v", err)
		return a.serveAuthPage(w, r)
	}

	a.AuthHeadersService.Process(r, claims)
	return false
}

func (a *Auth) serveViaFserver(w http.ResponseWriter, r *http.Request) bool {
	a.Log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("auth: via fserver")

	r.URL.Path = a.reURL.FindStringSubmatch(r.URL.Path)[1]
	a.fserver.ServeHTTP(w, r)

	return true
}

func (a *Auth) serveAuthPage(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusUnauthorized) // NOTE: 401 code helps js web-applications to determine expired tokens

	if err := a.page.Execute(w, a.Config); err != nil {
		a.Log.Errorf("auth: error executing auth page: %v", err)
	}

	return true
}

func (a *Auth) deleteCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     "/",
		HttpOnly: true,
		Domain:   a.Config.CookieDomain,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
}

func (a *Auth) updateCookie(w http.ResponseWriter, cookie *http.Cookie, value string) {
	cookie = &http.Cookie{
		Name:     cookie.Name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Domain:   a.Config.CookieDomain,
	}

	http.SetCookie(w, cookie)
}
