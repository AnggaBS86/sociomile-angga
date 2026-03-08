package service

import (
	"context"
	"strings"

	"sociomile-be/internal/domain/model"
	repository "sociomile-be/internal/domain/repository_interface"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, error) {
	// validate email and password input
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredential
	}

	// In production, compare hashed passwords.
	if user.Password != password {
		return nil, ErrInvalidCredential
	}

	return user, nil
}
