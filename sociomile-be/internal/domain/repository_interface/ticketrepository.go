package respository_interface

import (
	"context"
	"sociomile-be/internal/domain/model"
)

type TicketRepository interface {
	GetByConversationID(ctx context.Context, tenantID, conversationID int64) (*model.Ticket, error)
	Create(ctx context.Context, ticket *model.Ticket) error
	UpdateStatus(ctx context.Context, tenantID, ticketID int64, status string) error
	List(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, error)
}
