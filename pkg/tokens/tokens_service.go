package tokens

import (
	"sync"

	"github.com/goava/di"
	authapi "github.com/pridemon/outpost/pkg/auth_api"
	authheaders "github.com/pridemon/outpost/pkg/auth_headers"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/pridemon/outpost/pkg/models"
	"github.com/pridemon/outpost/pkg/repository"
	"github.com/pridemon/outpost/pkg/utils"
)

type TokensService struct {
	di.Inject

	AuthApi               *authapi.AuthApi
	AuthHeadersService    *authheaders.AuthHeadersService
	JwtService            *jwt.JwtService
	RefreshInfoRepository *repository.RefreshInfoRepository

	mutex            sync.Mutex
	processingTokens map[string]bool
}

func NewTokensService() *TokensService {
	return &TokensService{
		mutex:            sync.Mutex{},
		processingTokens: make(map[string]bool),
	}
}

func (srv *TokensService) IsProcessing(accessToken string) bool {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	return srv.processingTokens[accessToken]
}

func (srv *TokensService) ProcessAccessToken(accessToken string) (*jwt.JwtClaims, error) {
	srv.markAsProcessing(accessToken)
	defer srv.markAsFinished(accessToken)

	return srv.JwtService.CheckAccessToken(accessToken)
}

func (srv *TokensService) ProcessRefreshToken(accessToken string, refreshToken string) error {
	srv.markAsProcessing(accessToken)
	defer srv.markAsFinished(accessToken)

	return srv.RefreshInfoRepository.Insert(&models.RefreshInfo{
		Hash:         utils.GetMD5Hash(accessToken),
		RefreshToken: refreshToken,
	})
}

func (srv *TokensService) RefreshToken(accessToken string) (string, error) {
	srv.markAsProcessing(accessToken)
	defer srv.markAsFinished(accessToken)

	hash := utils.GetMD5Hash(accessToken)
	foundToken, err := srv.RefreshInfoRepository.Find(hash)
	if err != nil {
		return "", err
	}

	newTokens, err := srv.AuthApi.Refresh(foundToken.RefreshToken)
	if err != nil {
		return "", err
	}

	err = srv.RefreshInfoRepository.Insert(&models.RefreshInfo{
		Hash:         utils.GetMD5Hash(newTokens.AccessToken),
		RefreshToken: newTokens.RefreshToken,
	})
	if err != nil {
		return "", err
	}

	err = srv.RefreshInfoRepository.DeleteByHash(hash)
	if err != nil {
		return "", err
	}

	return newTokens.AccessToken, nil
}

func (srv *TokensService) markAsProcessing(accessToken string) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	srv.processingTokens[accessToken] = true
}

func (srv *TokensService) markAsFinished(accessToken string) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	delete(srv.processingTokens, accessToken)
}
