package service

import (
	"context"
	"errors"
	"time"

	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")

type AuthService interface {
	Login(ctx context.Context, username, password string) (*models.User, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
