package respository_interface

import (
	"context"
	"sociomile-be/internal/domain/model"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, activity *model.ActivityLog) error
}
