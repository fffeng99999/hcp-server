package handlers

import (
	"context"
	"time"

	pb "github.com/fffeng99999/hcp-server/api/generated/auth"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.svc.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if err == service.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if err == service.ErrInvalidCredentials {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, err
	}

	return &pb.LoginResponse{
		User: mapUserToProto(user),
	}, nil
}

func mapUserToProto(u *models.User) *pb.User {
	pbUser := &pb.User{
		Id:        u.ID.String(),
		Username:  u.Username,
		Role:      u.Role,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		Status:    u.Status,
	}
	if u.LastLogin != nil {
		pbUser.LastLogin = u.LastLogin.Format(time.RFC3339)
	}
	return pbUser
}
