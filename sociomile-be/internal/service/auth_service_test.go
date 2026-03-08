package service

import (
	"context"
	"errors"
	"testing"

	"sociomile-be/internal/domain/model"
)

func TestAuthServiceLogin(t *testing.T) {
	t.Run("returns invalid input for blank credentials", func(t *testing.T) {
		svc := NewAuthService(&fakeUserRepo{})

		_, err := svc.Login(context.Background(), " ", "")
		if !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns invalid credential when user not found", func(t *testing.T) {
		svc := NewAuthService(&fakeUserRepo{
			findByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return nil, nil
			},
		})

		_, err := svc.Login(context.Background(), "angga@email.com", "123456")
		if !errors.Is(err, ErrInvalidCredential) {
			t.Fatalf("expected ErrInvalidCredential, got %v", err)
		}
	})

	t.Run("returns invalid credential when password mismatch", func(t *testing.T) {
		svc := NewAuthService(&fakeUserRepo{
			findByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{Email: email, Password: "secret"}, nil
			},
		})

		_, err := svc.Login(context.Background(), "angga@email.com", "123456")
		if !errors.Is(err, ErrInvalidCredential) {
			t.Fatalf("expected ErrInvalidCredential, got %v", err)
		}
	})

	t.Run("returns user when credentials are valid", func(t *testing.T) {
		expected := &model.User{ID: 1, Email: "angga@email.com", Password: "123456"}
		svc := NewAuthService(&fakeUserRepo{
			findByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return expected, nil
			},
		})

		user, err := svc.Login(context.Background(), "angga@email.com", "123456")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		if user != expected {
			t.Fatalf("expected same user pointer returned")
		}
	})
}
