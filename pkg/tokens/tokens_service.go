package tokens

import (
	"errors"

	"github.com/goava/di"
	authapi "github.com/pridemon/outpost/pkg/auth_api"
	authheaders "github.com/pridemon/outpost/pkg/auth_headers"
	"github.com/pridemon/outpost/pkg/jwt"
	"github.com/pridemon/outpost/pkg/models"
	"github.com/pridemon/outpost/pkg/repository"
	"github.com/pridemon/outpost/pkg/utils"
	"golang.org/x/sync/semaphore"
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

	processingTokens map[string]*semaphore.Weighted
}

func NewTokensService() *TokensService {
	return &TokensService{
		processingTokens: make(map[string]*semaphore.Weighted),
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
	if err := srv.lockToken(accessToken); err != nil {
		return "", err
	}
	defer srv.unlockToken(accessToken)

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

func (srv *TokensService) lockToken(accessToken string) error {
	if srv.processingTokens[accessToken] == nil {
		srv.processingTokens[accessToken] = semaphore.NewWeighted(1)
	}

	if !srv.processingTokens[accessToken].TryAcquire(1) {
		return ErrBusyToken
	}

	return nil
}

func (srv *TokensService) unlockToken(accessToken string) {
	delete(srv.processingTokens, accessToken)
}
