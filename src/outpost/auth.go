package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

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
		a.rewriteAuthFormPost(r)

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

// This method replaces POST request from auth form to GET request
func (a *Auth) rewriteAuthFormPost(r *http.Request) {
	// NOTE: here we duplicate request body, because call to FormValue reads from r.Body buffer
	body, _ := ioutil.ReadAll(r.Body) // TODO: error check?
	buf1 := ioutil.NopCloser(bytes.NewBuffer(body))
	buf2 := ioutil.NopCloser(bytes.NewBuffer(body))

	r.Body = buf1 // first copy: for possible FormValue call

	// auth form sends special "__outpost__" input field
	if r.Method == "POST" && r.FormValue("__outpost__") == "__outpost__" {
		r.Method = "GET"
		r.Body = ioutil.NopCloser(strings.NewReader(""))
		r.ContentLength = 0
	} else {
		r.Body = buf2 // second copy: restore original body state
	}
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
