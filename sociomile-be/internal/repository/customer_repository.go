package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sociomile-be/internal/domain/model"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindByExternalID(ctx context.Context, tenantID int64, externalID string) (*model.Customer, error) {
	q := `SELECT id, tenant_id, external_id, created_at, updated_at FROM customers WHERE tenant_id = ? AND external_id = ? LIMIT 1`
	var c model.Customer

	err := r.db.QueryRowContext(ctx, q, tenantID, externalID).Scan(&c.ID, &c.TenantID, &c.ExternalID, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	res, err := r.db.ExecContext(ctx, `INSERT INTO customers (tenant_id, external_id) VALUES (?, ?)`, customer.TenantID, customer.ExternalID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	customer.ID = id

	return nil
}
