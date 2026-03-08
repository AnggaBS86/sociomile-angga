package asynqworker

import (
	"context"
	"fmt"
	"log"

	"sociomile-be/internal/domain/model"
	"sociomile-be/internal/service/events"

	"github.com/hibiken/asynq"
)

func (w *Worker) handleConversationEscalated(ctx context.Context, task *asynq.Task) error {
	event, err := decodeEvent(task)
	if err != nil {
		return err
	}

	tenantID := event.TenantID
	conversationID, err := getInt64(event.Payload, "conversation_id")
	if err != nil {
		return err
	}

	agentID, err := getInt64(event.Payload, "agent_id")
	if err != nil {
		return err
	}

	title := getString(event.Payload, "title", "Escalated conversation")
	description := getString(event.Payload, "description", "")
	priority := getString(event.Payload, "priority", "medium")

	conversation, err := w.conversationRepo.GetByID(ctx, tenantID, conversationID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return fmt.Errorf("conversation not found")
	}

	existing, err := w.ticketRepo.GetByConversationID(ctx, tenantID, conversationID)
	if err != nil {
		return err
	}
	if existing != nil {
		log.Printf("asynq skip create ticket: conversation=%d already has ticket=%d", conversationID, existing.ID)
		return nil
	}

	ticket := &model.Ticket{
		TenantID:       tenantID,
		ConversationID: conversationID,
		Title:          title,
		Description:    description,
		Status:         "open",
		Priority:       priority,
		AssignedAgentID: func() *int64 {
			if conversation.AssignedAgentID != nil {
				return conversation.AssignedAgentID
			}
			return &agentID
		}(),
	}

	if err := w.ticketRepo.Create(ctx, ticket); err != nil {
		return fmt.Errorf("create ticket: %w", err)
	}
	if err := w.listCache.InvalidateTicketList(ctx, tenantID); err != nil {
		return fmt.Errorf("invalidate ticket cache: %w", err)
	}

	if w.dispatcher != nil {
		err = w.dispatcher.Dispatch(ctx, model.DomainEvent{
			TenantID:   tenantID,
			EntityType: "ticket",
			EntityID:   ticket.ID,
			EventName:  events.EventTicketCreated,
			Payload: map[string]any{
				"ticket_id":       ticket.ID,
				"conversation_id": ticket.ConversationID,
				"priority":        ticket.Priority,
				"status":          ticket.Status,
			},
		})
		if err != nil {
			return err
		}
	}

	log.Printf("asynq processed event: %s conversation=%d ticket=%d", event.EventName, conversationID, ticket.ID)

	return nil
}
