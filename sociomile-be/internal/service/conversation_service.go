package service

import (
	"context"

	"sociomile-be/internal/domain/model"
	repository "sociomile-be/internal/domain/repository_interface"
	"sociomile-be/internal/service/cache"
	"sociomile-be/internal/service/events"
)

type ConversationService struct {
	customerRepo     repository.CustomerRepository
	conversationRepo repository.ConversationRepository
	listCache        cache.ListCache
	eventDispatcher  events.Dispatcher
}

type ConversationDetail struct {
	Conversation *model.Conversation `json:"conversation"`
	Messages     []model.Message     `json:"messages"`
}

func NewConversationService(customerRepo repository.CustomerRepository, conversationRepo repository.ConversationRepository, listCache cache.ListCache, eventDispatcher events.Dispatcher) *ConversationService {
	if listCache == nil {
		listCache = cache.NewNoopListCache()
	}

	if eventDispatcher == nil {
		eventDispatcher = events.NewNoopDispatcher()
	}

	return &ConversationService{customerRepo: customerRepo, conversationRepo: conversationRepo, listCache: listCache, eventDispatcher: eventDispatcher}
}

func (s *ConversationService) IngestChannelMessage(ctx context.Context, tenantID int64, customerExternalID, text string) (*model.Conversation, error) {
	if tenantID == 0 || customerExternalID == "" || text == "" {
		return nil, ErrInvalidInput
	}

	customer, err := s.customerRepo.FindByExternalID(ctx, tenantID, customerExternalID)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		customer = &model.Customer{TenantID: tenantID, ExternalID: customerExternalID}
		if err := s.customerRepo.Create(ctx, customer); err != nil {
			return nil, err
		}
	}

	conversation, err := s.conversationRepo.FindActiveByCustomerID(ctx, tenantID, customer.ID)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		conversation = &model.Conversation{TenantID: tenantID, CustomerID: customer.ID, Status: "open"}
		if err := s.conversationRepo.Create(ctx, conversation); err != nil {
			return nil, err
		}
	}

	msg := &model.Message{ConversationID: conversation.ID, SenderType: "customer", Message: text}
	if err := s.conversationRepo.AddMessage(ctx, msg); err != nil {
		return nil, err
	}

	if err := s.listCache.InvalidateConversationList(ctx, tenantID); err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *ConversationService) List(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, error) {
	if cached, found, err := s.listCache.GetConversationList(ctx, tenantID, filter); err == nil && found {
		return cached, nil
	} else if err != nil {
		return nil, err
	}

	conversations, err := s.conversationRepo.List(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}

	if err := s.listCache.SetConversationList(ctx, tenantID, filter, conversations); err != nil {
		return nil, err
	}

	return conversations, nil
}

func (s *ConversationService) GetDetail(ctx context.Context, tenantID, conversationID int64) (*ConversationDetail, error) {
	conversation, err := s.conversationRepo.GetByID(ctx, tenantID, conversationID)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, ErrNotFound
	}

	messages, err := s.conversationRepo.ListMessages(ctx, conversation.ID)
	if err != nil {
		return nil, err
	}

	return &ConversationDetail{Conversation: conversation, Messages: messages}, nil
}

func (s *ConversationService) AgentReply(ctx context.Context, tenantID, conversationID, agentID int64, role, text string) error {
	if text == "" {
		return ErrInvalidInput
	}

	conversation, err := s.conversationRepo.GetByID(ctx, tenantID, conversationID)
	if err != nil {
		return err
	}

	if conversation == nil {
		return ErrNotFound
	}

	if role == model.RoleAgent && conversation.AssignedAgentID != nil && *conversation.AssignedAgentID != agentID {
		return ErrForbidden
	}

	if conversation.AssignedAgentID == nil {
		err = s.eventDispatcher.Dispatch(ctx, model.DomainEvent{
			TenantID:   tenantID,
			EntityType: "conversation",
			EntityID:   conversation.ID,
			EventName:  events.EventConversationAssigned,
			Payload: map[string]any{
				"conversation_id": conversation.ID,
				"agent_id":        agentID,
				"status":          "assigned",
			},
		})

		if err != nil {
			return err
		}
	}

	msg := &model.Message{ConversationID: conversation.ID, SenderType: "agent", Message: text}
	if err := s.conversationRepo.AddMessage(ctx, msg); err != nil {
		return err
	}

	return s.listCache.InvalidateConversationList(ctx, tenantID)
}
