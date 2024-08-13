package usecases

import (
	"time"

	"github.com/emaforlin/auth-service/internal/config"
	"github.com/emaforlin/auth-service/pkg/entities"
	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenUsecase interface {
	GetClaims(token string) *entities.CustomClaims
	NewToken(role string) (string, error)
}

type tokenUsecase struct {
	cfg config.Config
}

// NewToken implements TokenUsecase.
func (t *tokenUsecase) NewToken(role string) (string, error) {
	claims := &entities.CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.cfg.Jwt.Ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(t.cfg.Jwt.Secret)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// GetClaims implements TokenUsecase.
func (t *tokenUsecase) GetClaims(token string) *entities.CustomClaims {
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
