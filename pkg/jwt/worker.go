package jwt

import (
	"time"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/repository"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	di.Inject

	Log              *logrus.Logger
	Config           *JwtConfig
	TokensRepository *repository.TokensRepository
}

func (w *Worker) Run() {
	for {
		checkTime := time.Now().Add(-1 * w.Config.RefreshTokenTTL) // find edge time for created_at

		n, err := w.TokensRepository.DeleteWithCreatedAtLT(checkTime)
		if err != nil {
			w.Log.Errorf("jwt.worker: error during sql query: %v", err)
		} else {
			w.Log.Printf("jwt.worker: deleted %d records", n)
		}

		time.Sleep(w.Config.WorkerDelay)
	}
}
