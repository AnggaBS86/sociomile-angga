package service

import (
	"context"
	"errors"
	"testing"

	"sociomile-be/internal/domain/model"
	"sociomile-be/internal/service/events"
)

func TestTicketServiceEscalateDispatchesAsyncEvent(t *testing.T) {
	dispatcher := &recordingDispatcher{}
	svc := NewTicketService(
		&fakeConversationRepo{
			getByIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
				return &model.Conversation{ID: conversationID, TenantID: tenantID}, nil
			},
		},
		&fakeTicketRepo{
			getByConversationIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Ticket, error) {
				return nil, nil
			},
		},
		&fakeListCache{},
		dispatcher,
	)

	ticket, err := svc.Escalate(context.Background(), 1, 44, 8, "Escalated issue", "Need internal follow up", "high")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ticket.Status != "queued" {
		t.Fatalf("expected queued ticket status, got %s", ticket.Status)
	}

	if len(dispatcher.events) != 1 {
		t.Fatalf("expected one event dispatched, got %d", len(dispatcher.events))
	}

	if dispatcher.events[0].EventName != events.EventConversationEscalated {
		t.Fatalf("expected %s, got %s", events.EventConversationEscalated, dispatcher.events[0].EventName)
	}
}

func TestTicketServiceEscalateRejectsExistingTicket(t *testing.T) {
	svc := NewTicketService(
		&fakeConversationRepo{
			getByIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
				return &model.Conversation{ID: conversationID, TenantID: tenantID}, nil
			},
		},
		&fakeTicketRepo{
			getByConversationIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Ticket, error) {
				return &model.Ticket{ID: 99, ConversationID: conversationID}, nil
			},
		},
		&fakeListCache{},
		&recordingDispatcher{},
	)

	_, err := svc.Escalate(context.Background(), 1, 44, 8, "title", "desc", "high")
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestTicketServiceUpdateStatusRejectsUnsupportedStatus(t *testing.T) {
	svc := NewTicketService(
		&fakeConversationRepo{},
		&fakeTicketRepo{},
		&fakeListCache{},
		&recordingDispatcher{},
	)

	err := svc.UpdateStatus(context.Background(), 1, 10, "resolved")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
