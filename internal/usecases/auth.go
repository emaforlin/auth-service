package usecases

import (
	"github.com/emaforlin/auth-service/internal/config"
	pb "github.com/emaforlin/auth-service/pkg/pb/protos"
)

type AuthUsecase interface {
	CheckPermissionScope(*pb.AuthorizationRequest) []string
}

type authUsecase struct {
	cfg          config.Config
	tokenManager TokenUsecase
}

// CheckPermissionScope implements AuthUsecase.
func (u *authUsecase) CheckPermissionScope(in *pb.AuthorizationRequest) []string {
	claims := u.tokenManager.Decode(in.GetToken())
	permissions := u.cfg.AccessControl[claims.Role]
	return permissions
}

func NewAuthUsecase(c *config.Config) AuthUsecase {
	return &authUsecase{
		cfg:          *c,
		tokenManager: NewTokenUsecase(c),
	}
}
