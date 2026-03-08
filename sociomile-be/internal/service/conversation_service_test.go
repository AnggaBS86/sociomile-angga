package service

import (
	"context"
	"errors"
	"testing"

	"sociomile-be/internal/domain/model"
	"sociomile-be/internal/service/events"
)

func TestConversationServiceIngestChannelMessageCreatesConversationAndMessage(t *testing.T) {
	var createdConversation *model.Conversation
	var createdCustomer *model.Customer
	var addedMessage *model.Message

	svc := NewConversationService(
		&fakeCustomerRepo{
			findByExternalIDFn: func(ctx context.Context, tenantID int64, externalID string) (*model.Customer, error) {
				return nil, nil
			},
			createFn: func(ctx context.Context, customer *model.Customer) error {
				customer.ID = 101
				createdCustomer = customer
				return nil
			},
		},
		&fakeConversationRepo{
			findActiveByCustomerIDFn: func(ctx context.Context, tenantID, customerID int64) (*model.Conversation, error) {
				return nil, nil
			},
			createFn: func(ctx context.Context, conversation *model.Conversation) error {
				conversation.ID = 202
				createdConversation = conversation
				return nil
			},
			addMessageFn: func(ctx context.Context, message *model.Message) error {
				addedMessage = message
				return nil
			},
		},
		&fakeListCache{},
		&recordingDispatcher{},
	)

	conversation, err := svc.IngestChannelMessage(context.Background(), 1, "cust-1", "hello")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if createdCustomer == nil || createdCustomer.TenantID != 1 {
		t.Fatalf("expected customer to be created for tenant 1")
	}

	if createdConversation == nil || createdConversation.CustomerID != 101 || createdConversation.Status != "open" {
		t.Fatalf("expected conversation to be created with open status")
	}

	if addedMessage == nil || addedMessage.ConversationID != 202 || addedMessage.SenderType != "customer" {
		t.Fatalf("expected customer message to be added")
	}

	if conversation.ID != 202 {
		t.Fatalf("expected returned conversation id 202, got %d", conversation.ID)
	}
}

func TestConversationServiceAgentReplyDispatchesAssignmentWhenUnassigned(t *testing.T) {
	dispatcher := &recordingDispatcher{}
	var addedMessage *model.Message

	svc := NewConversationService(
		&fakeCustomerRepo{},
		&fakeConversationRepo{
			getByIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
				return &model.Conversation{ID: conversationID, TenantID: tenantID, Status: "open"}, nil
			},
			addMessageFn: func(ctx context.Context, message *model.Message) error {
				addedMessage = message
				return nil
			},
		},
		&fakeListCache{},
		dispatcher,
	)

	err := svc.AgentReply(context.Background(), 1, 55, 7, model.RoleAgent, "reply")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(dispatcher.events) != 1 {
		t.Fatalf("expected 1 event dispatched, got %d", len(dispatcher.events))
	}

	event := dispatcher.events[0]
	if event.EventName != events.EventConversationAssigned {
		t.Fatalf("expected %s, got %s", events.EventConversationAssigned, event.EventName)
	}

	if addedMessage == nil || addedMessage.SenderType != "agent" {
		t.Fatalf("expected agent message to be added")
	}
}

func TestConversationServiceAgentReplyForbidsDifferentAssignedAgent(t *testing.T) {
	assignedAgentID := int64(99)
	svc := NewConversationService(
		&fakeCustomerRepo{},
		&fakeConversationRepo{
			getByIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
				return &model.Conversation{ID: conversationID, TenantID: tenantID, AssignedAgentID: &assignedAgentID}, nil
			},
		},
		&fakeListCache{},
		&recordingDispatcher{},
	)

	err := svc.AgentReply(context.Background(), 1, 55, 7, model.RoleAgent, "reply")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestConversationServiceGetDetailReturnsConversationAndMessages(t *testing.T) {
	svc := NewConversationService(
		&fakeCustomerRepo{},
		&fakeConversationRepo{
			getByIDFn: func(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
				return &model.Conversation{ID: conversationID, TenantID: tenantID}, nil
			},
			listMessagesFn: func(ctx context.Context, conversationID int64) ([]model.Message, error) {
				return []model.Message{{ID: 1, ConversationID: conversationID, Message: "hello"}}, nil
			},
		},
		&fakeListCache{},
		&recordingDispatcher{},
	)

	detail, err := svc.GetDetail(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if detail.Conversation.ID != 10 || len(detail.Messages) != 1 {
		t.Fatalf("expected conversation detail with one message")
	}
}
