package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goava/di"
	"github.com/pridemon/outpost/internal"
	"github.com/pridemon/outpost/pkg/auth"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/pridemon/outpost/pkg/proxy"
	"github.com/pridemon/outpost/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	di.SetTracer(&di.StdTracer{})

	container, err := di.New(
		internal.HttpModule,
		internal.ViperModule,
		internal.LogrusModule,
		internal.RestyModule,

		internal.JwtModule,
		internal.AuthModule,
		internal.AuthHeadersModule,
		internal.AuthApiModule,
		internal.ProxyModule,
		internal.SqlModule,
		internal.TokensModule,

		di.Invoke(RunJwtWorker),
		di.Invoke(RunServer),
	)
	if err != nil {
		log.Fatalln("error:", err)
	}

	container.Cleanup()
}

func RunJwtWorker(w *jwt.Worker) {
	go w.Run()
}

func RunServer(auth *auth.Auth, proxy *proxy.Proxy, httpConfig *internal.HttpConfig, log *logrus.Logger) {
	utils.DisableSSLVerification()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.PreventIndexing(w)

		if auth.TryServeHTTP(w, r) {
			return
		}

		proxy.TryServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpConfig.Port), nil))
}
