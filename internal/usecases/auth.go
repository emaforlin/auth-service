package usecases

import (
	"time"

	accountsPb "github.com/emaforlin/accounts-service/x/handlers/grpc/protos"
	"github.com/emaforlin/auth-service/internal/config"
	pb "github.com/emaforlin/auth-service/pkg/pb/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthUsecase interface {
	Login(*pb.LoginRequest) string
	CheckPermissionScope(*pb.AuthorizationRequest) []string
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
