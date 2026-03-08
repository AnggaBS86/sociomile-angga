package respository_interface

import (
	"context"
	"sociomile-be/internal/domain/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}
