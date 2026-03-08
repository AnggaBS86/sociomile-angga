package respository_interface

import (
	"context"
	"sociomile-be/internal/domain/model"
)

type CustomerRepository interface {
	FindByExternalID(ctx context.Context, tenantID int64, externalID string) (*model.Customer, error)
	Create(ctx context.Context, customer *model.Customer) error
}
