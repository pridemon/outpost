package internal

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/auth"
	"github.com/spf13/viper"
)

var AuthModule = di.Options(
	di.Provide(AuthConfigProvider),
	di.Provide(AuthProvider),
)

func AuthConfigProvider(v *viper.Viper) (*auth.AuthConfig, error) {
	var config auth.AuthConfig
	err := v.UnmarshalKey("auth", &config)
	return &config, err
}

func AuthProvider(config *auth.AuthConfig) (*auth.Auth, error) {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return nil, err
	}
	dir = dir + "/html"

	fserver := http.FileServer(http.Dir(dir))

	bytes, err := ioutil.ReadFile(dir + "/index.html")
	if err != nil {
		return nil, err
	}

	page, err := template.New("auth-page").Parse(string(bytes))
	if err != nil {
		return nil, err
	}

	return auth.NewAuth(
		fserver,
		regexp.MustCompile("^/__outpost__/(.*)$"),
		page,
	), nil
}
