package main

import (
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type Auth struct {
	Title string
	Icon  string

	fserver http.Handler
	reURL   *regexp.Regexp
	page    *template.Template
	hash    string
}

func (a *Auth) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if a.checkCookie(r) {
		return false
	}

	log.WithField("method", r.Method).WithField("url", r.URL.String()).Debug("via fserver")

	if !a.reURL.MatchString(r.URL.Path) {
		a.serveAuthPage(w, r)
	} else {
		r.URL.Path = a.reURL.FindStringSubmatch(r.URL.Path)[1]
		a.fserver.ServeHTTP(w, r)
	}

	return true
}

func (a *Auth) checkCookie(r *http.Request) bool {
	cookie, err := r.Cookie("outpost")
	if err != nil {
		return false
	}

	return cookie.Value == a.hash
}

func (a *Auth) SetLoginPassw(login, passw string) {
	h := sha1.New()
	h.Write([]byte(login + "|" + passw))
	a.hash = hex.EncodeToString(h.Sum(nil))
}

func (a *Auth) serveAuthPage(w http.ResponseWriter, r *http.Request) {
	if err := a.page.Execute(w, a); err != nil {
		log.Error(err)
	}
}
