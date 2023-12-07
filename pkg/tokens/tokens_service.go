package tokens

import (
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

	AuthApi            *authapi.AuthApi
	AuthHeadersService *authheaders.AuthHeadersService
	JwtService         *jwt.JwtService
	TokensRepository   *repository.TokensRepository
}

func (srv *TokensService) ProcessTokens(accessToken string, refreshToken string) (*jwt.JwtClaims, error) {
	claims, err := srv.ProcessAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	if refreshToken != "" {
		err = srv.TokensRepository.Insert(&models.Token{
			Hash:         utils.GetMD5Hash(accessToken),
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func (srv *TokensService) ProcessAccessToken(accessToken string) (*jwt.JwtClaims, error) {
	return srv.JwtService.CheckAccessToken(accessToken)
}

func (srv *TokensService) RefreshToken(accessToken string) (string, error) {
	hash := utils.GetMD5Hash(accessToken)
	foundToken, err := srv.TokensRepository.FindToken(hash)
	if err != nil {
		return "", err
	}

	newTokens, err := srv.AuthApi.Refresh(foundToken.RefreshToken)
	if err != nil {
		return "", err
	}

	err = srv.TokensRepository.Insert(&models.Token{
		Hash:         utils.GetMD5Hash(newTokens.AccessToken),
		RefreshToken: newTokens.RefreshToken,
	})
	if err != nil {
		return "", err
	}

	err = srv.TokensRepository.DeleteByHash(hash)
	if err != nil {
		return "", err
	}

	return newTokens.AccessToken, nil
}
