package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose logging.").Short('v').Bool()
)

func main() {
	kingpin.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	DisableSSLVerification()

	auth := ConfiguredAuth()
	proxy := ConfiguredProxy()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if auth.TryServeHTTP(w, r) {
			return
		}

		proxy.TryServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
