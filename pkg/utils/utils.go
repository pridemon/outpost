package utils

import (
	"crypto/tls"
	"net/http"
)

func PreventIndexing(w http.ResponseWriter) {
	w.Header().Set("X-Robots-Tag", "noindex, nofollow, nosnippet, noarchive")
}

func DisableSSLVerification() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
