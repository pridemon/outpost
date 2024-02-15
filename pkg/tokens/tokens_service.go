package tokens

import (
	"errors"
	"sync"

	"github.com/goava/di"
	authapi "github.com/pridemon/outpost/pkg/auth_api"
	authheaders "github.com/pridemon/outpost/pkg/auth_headers"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/pridemon/outpost/pkg/models"
	"github.com/pridemon/outpost/pkg/repository"
	"github.com/pridemon/outpost/pkg/utils"
)

var (
	ErrBusyToken = errors.New("token is busy")
)

type TokensService struct {
	di.Inject

	AuthApi               *authapi.AuthApi
	AuthHeadersService    *authheaders.AuthHeadersService
	JwtService            *jwt.JwtService
	RefreshInfoRepository *repository.RefreshInfoRepository

	processingTokens map[string]*sync.WaitGroup
	mutex            sync.Mutex
}

func NewTokensService() *TokensService {
	return &TokensService{
		processingTokens: make(map[string]*sync.WaitGroup),
	}
}

func (srv *TokensService) ProcessAccessToken(accessToken string) (*jwt.JwtClaims, error) {
	return srv.JwtService.CheckAccessToken(accessToken)
}

func (srv *TokensService) ProcessRefreshToken(accessToken string, refreshToken string) error {
	return srv.RefreshInfoRepository.Insert(&models.RefreshInfo{
		Hash:         utils.GetMD5Hash(accessToken),
		RefreshToken: refreshToken,
	})
}

func (srv *TokensService) RefreshToken(accessToken string) (string, error) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	if srv.processingTokens[accessToken] == nil {
		srv.processingTokens[accessToken] = &sync.WaitGroup{}
		srv.processingTokens[accessToken].Add(1)

		newToken, err := srv.refreshToken(accessToken)

		srv.processingTokens[accessToken].Done()
		return newToken, err
	}

	srv.processingTokens[accessToken].Wait()
	delete(srv.processingTokens, accessToken)
	return "", ErrBusyToken
}

func (srv *TokensService) refreshToken(accessToken string) (string, error) {
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
