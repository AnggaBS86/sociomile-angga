package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sociomile-be/internal/domain/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	q := `SELECT id, tenant_id, email, password, role, created_at, updated_at FROM users WHERE email = ? LIMIT 1`
	var user model.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
