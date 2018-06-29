package main

import "net/http"

func PreventIndexing(w http.ResponseWriter) {
	w.Header().Set("X-Robots-Tag", "noindex, nofollow, nosnippet, noarchive")
}
