package main

import (
	"crypto/tls"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func ConfiguredAuth() *Auth {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = dir + "/../html"

	log.Info(dir)
	fserver := http.FileServer(http.Dir(dir))

	bytes, err := ioutil.ReadFile(dir + "/index.html")
	if err != nil {
		log.Fatal(err)
	}

	page, err := template.New("auth-page").Parse(string(bytes))
	if err != nil {
		log.Fatal(err)
	}

	auth := &Auth{
		Title: GetenvStr("OUTPOST_TITLE", "Outpost"),
		Icon:  GetenvStr("OUTPOST_ICON", ""), //"https://is3-ssl.mzstatic.com/image/thumb/Purple118/v4/c1/3d/a5/c13da51a-5152-29d3-b668-4547e8873cc6/mzl.nhnzrmvu.png/230x0w.jpg",

		fserver: fserver,
		reURL:   regexp.MustCompile("^/__outpost__/(.*)$"),
		page:    page,
	}

	auth.SetLoginPassw(
		GetenvStrFatal("OUTPOST_LOGIN"),
		GetenvStrFatal("OUTPOST_PASSW"),
	)
	return auth
}

func DisableSSLVerification() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func ConfiguredProxy() *Proxy {

	proxy := NewProxy()

	if target, err := url.Parse(GetenvStrFatal("OUTPOST_TARGET")); err == nil {
		proxy.Target = target
	} else {
		log.Fatal("OUTPOST_TARGET is not valid url")
	}

	proxy.Host = GetenvStr("OUTPOST_HOST", "")

	return proxy
}

func GetenvStr(name string, fallback string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	return fallback
}

func GetenvStrFatal(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	log.Fatalf("Can't find env var '%s'", name)
	return ""
}
