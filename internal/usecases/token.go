package usecases

import (
	"github.com/emaforlin/auth-service/internal/config"
	"github.com/emaforlin/auth-service/pkg/entities"
	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenUsecase interface {
	Decode(token string) *entities.CustomClaims
}

type tokenUsecase struct {
	cfg config.Config
}

// GetClaims implements TokenUsecase.
func (t *tokenUsecase) Decode(token string) *entities.CustomClaims {
	jwt, err := jwt.ParseWithClaims(token, &entities.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return t.cfg.Jwt.Secret, nil
	})
	if err != nil {
		return nil
	} else if claims, ok := jwt.Claims.(*entities.CustomClaims); ok {
		return claims
	}
	return nil
}

func NewTokenUsecase(c *config.Config) TokenUsecase {
	return &tokenUsecase{
		cfg: *c,
	}
}
