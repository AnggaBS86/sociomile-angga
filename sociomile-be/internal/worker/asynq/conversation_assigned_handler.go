package asynqworker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func (w *Worker) handleConversationAssigned(ctx context.Context, task *asynq.Task) error {
	event, err := decodeEvent(task)
	if err != nil {
		return err
	}

	conversationID, err := getInt64(event.Payload, "conversation_id")
	if err != nil {
		return err
	}

	agentID, err := getInt64(event.Payload, "agent_id")
	if err != nil {
		return err
	}

	status := getString(event.Payload, "status", "assigned")

	if err := w.conversationRepo.UpdateAssignment(ctx, conversationID, agentID, status); err != nil {
		return fmt.Errorf("update assignment: %w", err)
	}
	if err := w.listCache.InvalidateConversationList(ctx, event.TenantID); err != nil {
		return fmt.Errorf("invalidate conversation cache: %w", err)
	}

	log.Printf("asynq processed event: %s conversation=%d agent=%d", event.EventName, conversationID, agentID)

	return nil
}
