package service

import (
	"context"
	"strings"

	"sociomile-be/internal/domain/model"
	repository "sociomile-be/internal/domain/repository_interface"
	"sociomile-be/internal/service/cache"
	"sociomile-be/internal/service/events"
)

const (
	CONVERSATION_STATUS_OPEN     string = "open"
	CONVERSATION_STATUS_ASSIGNED string = "assigned"
	CONVERSATION_STATUS_CLOSED   string = "closed"
)

type TicketService struct {
	conversationRepo repository.ConversationRepository
	ticketRepo       repository.TicketRepository
	listCache        cache.ListCache
	eventDispatcher  events.Dispatcher
}

func NewTicketService(conversationRepo repository.ConversationRepository, ticketRepo repository.TicketRepository, listCache cache.ListCache, eventDispatcher events.Dispatcher) *TicketService {
	if listCache == nil {
		listCache = cache.NewNoopListCache()
	}
	if eventDispatcher == nil {
		eventDispatcher = events.NewNoopDispatcher()
	}
	return &TicketService{conversationRepo: conversationRepo, ticketRepo: ticketRepo, listCache: listCache, eventDispatcher: eventDispatcher}
}

func (s *TicketService) Escalate(ctx context.Context, tenantID, conversationID, agentID int64, title, description, priority string) (*model.Ticket, error) {
	conversation, err := s.conversationRepo.GetByID(ctx, tenantID, conversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, ErrNotFound
	}

	existing, err := s.ticketRepo.GetByConversationID(ctx, tenantID, conversationID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyExists
	}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if title == "" {
		title = "Escalated conversation"
	}
	if priority == "" {
		priority = "medium"
	}

	// Async command: worker will create ticket and emit ticket.created event.
	if err := s.eventDispatcher.Dispatch(ctx, model.DomainEvent{
		TenantID:   tenantID,
		EntityType: "conversation",
		EntityID:   conversationID,
		EventName:  events.EventConversationEscalated,
		Payload: map[string]any{
			"conversation_id": conversationID,
			"agent_id":        agentID,
			"title":           title,
			"description":     description,
			"priority":        priority,
		},
	}); err != nil {
		return nil, err
	}

	return &model.Ticket{
		TenantID:       tenantID,
		ConversationID: conversationID,
		Title:          title,
		Description:    description,
		Status:         "queued",
		Priority:       priority,
	}, nil
}

func (s *TicketService) UpdateStatus(ctx context.Context, tenantID, ticketID int64, status string) error {
	status = strings.TrimSpace(status)

	switch status {
	case CONVERSATION_STATUS_OPEN, CONVERSATION_STATUS_ASSIGNED, CONVERSATION_STATUS_CLOSED:
	default:
		return ErrInvalidInput
	}
	if err := s.ticketRepo.UpdateStatus(ctx, tenantID, ticketID, status); err != nil {
		return err
	}

	return s.listCache.InvalidateTicketList(ctx, tenantID)
}

func (s *TicketService) List(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, error) {
	if cached, found, err := s.listCache.GetTicketList(ctx, tenantID, filter); err == nil && found {
		return cached, nil
	} else if err != nil {
		return nil, err
	}

	tickets, err := s.ticketRepo.List(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}
	if err := s.listCache.SetTicketList(ctx, tenantID, filter, tickets); err != nil {
		return nil, err
	}

	return tickets, nil
}
