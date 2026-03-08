package repositories

import (
	"context"
	"database/sql"

	"sociomile-be/internal/domain/model"
)

type ActivityLogRepository struct {
	db *sql.DB
}

func NewActivityLogRepository(db *sql.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Create(ctx context.Context, activity *model.ActivityLog) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO activity_logs (tenant_id, entity_type, entity_id, event_name, payload) VALUES (?, ?, ?, ?, ?)`,
		activity.TenantID,
		activity.EntityType,
		activity.EntityID,
		activity.EventName,
		activity.Payload,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	activity.ID = id
	return nil
}
