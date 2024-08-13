package usecases

import (
	"time"

	accountsPb "github.com/emaforlin/accounts-service/x/handlers/grpc/protos"
	"github.com/emaforlin/auth-service/internal/config"
	pb "github.com/emaforlin/auth-service/pkg/pb/protos"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthUsecase interface {
	Login(*pb.LoginRequest) string
	IsAuthorized(*pb.AuthorizationRequest) bool
}

type authUsecase struct {
	cfg          config.Config
	tokenManager TokenUsecase
}

// Login implements AuthUsecase.
func (u *authUsecase) Login(in *pb.LoginRequest) string {
	// cl, err := grpc.NewClient(u.cfg.Dependencies["accounts"])
	cl, err := grpc.NewClient("localhost:50014", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("error connecting to external service")
	}

	accClient := accountsPb.NewAccountsClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	checkData := new(accountsPb.CheckUserPassRequest)

	switch in.Identifiers.(type) {
	case *pb.LoginRequest_Email:
		checkData.Identifiers = &accountsPb.CheckUserPassRequest_Email{Email: in.GetEmail()}
	case *pb.LoginRequest_Username:
		checkData.Identifiers = &accountsPb.CheckUserPassRequest_Username{Username: in.GetUsername()}
	case *pb.LoginRequest_PhoneNumber:
		checkData.Identifiers = &accountsPb.CheckUserPassRequest_PhoneNumber{PhoneNumber: in.GetPhoneNumber()}
	}

	checkData.Password = in.GetPassword()
	checkData.Role = in.GetRole()

	valid, err := accClient.CheckLoginData(ctx, checkData)
	if !valid.GetOk() || err != nil {
		return ""
	}

	token, err := u.tokenManager.NewToken(in.GetRole())
	if err != nil {
		return ""
	}
	return token
}

// CheckPermissionScope implements AuthUsecase.
func (u *authUsecase) IsAuthorized(in *pb.AuthorizationRequest) bool {
	tokenStr := in.GetToken()

	claims := u.tokenManager.GetClaims(tokenStr)
	permissions := u.cfg.AccessControl[claims.Role]

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return u.cfg.Jwt.Secret, nil
	})
	if err != nil {
		return false
	}
	if !token.Valid {
		return false
	}
}

func NewAuthUsecase(c *config.Config) AuthUsecase {
	return &authUsecase{
		cfg:          *c,
		tokenManager: NewTokenUsecase(c),
	}
}
