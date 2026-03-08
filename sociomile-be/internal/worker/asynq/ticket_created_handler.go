package asynqworker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func (w *Worker) handleTicketCreated(ctx context.Context, task *asynq.Task) error {
	event, err := decodeEvent(task)
	if err != nil {
		return err
	}
	if err := w.listCache.InvalidateTicketList(ctx, event.TenantID); err != nil {
		return fmt.Errorf("invalidate ticket cache: %w", err)
	}

	log.Printf("asynq processed event: %s tenant_id=%d entity=%s:%d",
		event.EventName, event.TenantID, event.EntityType, event.EntityID)

	return nil
}
