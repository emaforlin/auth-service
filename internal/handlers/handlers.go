package handlers

import (
	"context"
	"errors"

	"github.com/emaforlin/auth-service/internal/usecases"
	pb "github.com/emaforlin/auth-service/pkg/pb/protos"
	hclog "github.com/hashicorp/go-hclog"
)

type authServerImpl struct {
	pb.UnimplementedAuthServer
	log     hclog.Logger
	usecase usecases.AuthUsecase
}

func (h *authServerImpl) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	h.log.Info("Handle Login")
	token := h.usecase.Login(in)
	if token == "" {
		return nil, errors.New("authentication failed")
	}
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (h *authServerImpl) CheckPermissionScope(ctx context.Context, in *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	h.log.Info("Handle CheckPermissionScope")
	return &pb.AuthorizationResponse{AllowedMethods: h.usecase.CheckPermissionScope(in)}, nil
}

func NewAuthHandler(l hclog.Logger, u usecases.AuthUsecase) *authServerImpl {
	return &authServerImpl{
		log:     l,
		usecase: u,
	}
}
