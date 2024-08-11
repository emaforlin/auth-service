package handlers

import (
	"context"

	"github.com/emaforlin/auth-service/internal/usecases"
	pb "github.com/emaforlin/auth-service/pkg/pb/protos"
	hclog "github.com/hashicorp/go-hclog"
)

type authServerImpl struct {
	pb.UnimplementedAuthServer
	log     hclog.Logger
	usecase usecases.AuthUsecase
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
