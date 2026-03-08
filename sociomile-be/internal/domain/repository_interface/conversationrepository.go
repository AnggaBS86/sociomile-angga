package respository_interface

import (
	"context"
	"sociomile-be/internal/domain/model"
)

type ConversationRepository interface {
	FindActiveByCustomerID(ctx context.Context, tenantID, customerID int64) (*model.Conversation, error)
	Create(ctx context.Context, conversation *model.Conversation) error
	AddMessage(ctx context.Context, message *model.Message) error
	GetByID(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error)
	List(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, error)
	ListMessages(ctx context.Context, conversationID int64) ([]model.Message, error)
	UpdateAssignment(ctx context.Context, conversationID int64, assignedAgentID int64, status string) error
}
