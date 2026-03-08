package service

import (
	"context"

	"sociomile-be/internal/domain/model"
)

type fakeUserRepo struct {
	findByEmailFn func(ctx context.Context, email string) (*model.User, error)
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if f.findByEmailFn != nil {
		return f.findByEmailFn(ctx, email)
	}
	return nil, nil
}

type fakeCustomerRepo struct {
	findByExternalIDFn func(ctx context.Context, tenantID int64, externalID string) (*model.Customer, error)
	createFn           func(ctx context.Context, customer *model.Customer) error
}

func (f *fakeCustomerRepo) FindByExternalID(ctx context.Context, tenantID int64, externalID string) (*model.Customer, error) {
	if f.findByExternalIDFn != nil {
		return f.findByExternalIDFn(ctx, tenantID, externalID)
	}
	return nil, nil
}

func (f *fakeCustomerRepo) Create(ctx context.Context, customer *model.Customer) error {
	if f.createFn != nil {
		return f.createFn(ctx, customer)
	}
	return nil
}

type fakeConversationRepo struct {
	findActiveByCustomerIDFn func(ctx context.Context, tenantID, customerID int64) (*model.Conversation, error)
	createFn                 func(ctx context.Context, conversation *model.Conversation) error
	addMessageFn             func(ctx context.Context, message *model.Message) error
	getByIDFn                func(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error)
	listFn                   func(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, error)
	listMessagesFn           func(ctx context.Context, conversationID int64) ([]model.Message, error)
	updateAssignmentFn       func(ctx context.Context, conversationID int64, assignedAgentID int64, status string) error
}

func (f *fakeConversationRepo) FindActiveByCustomerID(ctx context.Context, tenantID, customerID int64) (*model.Conversation, error) {
	if f.findActiveByCustomerIDFn != nil {
		return f.findActiveByCustomerIDFn(ctx, tenantID, customerID)
	}
	return nil, nil
}

func (f *fakeConversationRepo) Create(ctx context.Context, conversation *model.Conversation) error {
	if f.createFn != nil {
		return f.createFn(ctx, conversation)
	}
	return nil
}

func (f *fakeConversationRepo) AddMessage(ctx context.Context, message *model.Message) error {
	if f.addMessageFn != nil {
		return f.addMessageFn(ctx, message)
	}
	return nil
}

func (f *fakeConversationRepo) GetByID(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, tenantID, conversationID)
	}
	return nil, nil
}

func (f *fakeConversationRepo) List(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, error) {
	if f.listFn != nil {
		return f.listFn(ctx, tenantID, filter)
	}
	return nil, nil
}

func (f *fakeConversationRepo) ListMessages(ctx context.Context, conversationID int64) ([]model.Message, error) {
	if f.listMessagesFn != nil {
		return f.listMessagesFn(ctx, conversationID)
	}
	return nil, nil
}

func (f *fakeConversationRepo) UpdateAssignment(ctx context.Context, conversationID int64, assignedAgentID int64, status string) error {
	if f.updateAssignmentFn != nil {
		return f.updateAssignmentFn(ctx, conversationID, assignedAgentID, status)
	}
	return nil
}

type fakeTicketRepo struct {
	getByConversationIDFn func(ctx context.Context, tenantID, conversationID int64) (*model.Ticket, error)
	createFn              func(ctx context.Context, ticket *model.Ticket) error
	updateStatusFn        func(ctx context.Context, tenantID, ticketID int64, status string) error
	listFn                func(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, error)
}

func (f *fakeTicketRepo) GetByConversationID(ctx context.Context, tenantID, conversationID int64) (*model.Ticket, error) {
	if f.getByConversationIDFn != nil {
		return f.getByConversationIDFn(ctx, tenantID, conversationID)
	}
	return nil, nil
}

func (f *fakeTicketRepo) Create(ctx context.Context, ticket *model.Ticket) error {
	if f.createFn != nil {
		return f.createFn(ctx, ticket)
	}
	return nil
}

func (f *fakeTicketRepo) UpdateStatus(ctx context.Context, tenantID, ticketID int64, status string) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, tenantID, ticketID, status)
	}
	return nil
}

func (f *fakeTicketRepo) List(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, error) {
	if f.listFn != nil {
		return f.listFn(ctx, tenantID, filter)
	}
	return nil, nil
}

type fakeListCache struct {
	getConversationListFn    func(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, bool, error)
	setConversationListFn    func(ctx context.Context, tenantID int64, filter model.ConversationFilter, value []model.Conversation) error
	invalidateConversationFn func(ctx context.Context, tenantID int64) error
	getTicketListFn          func(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, bool, error)
	setTicketListFn          func(ctx context.Context, tenantID int64, filter model.TicketFilter, value []model.Ticket) error
	invalidateTicketFn       func(ctx context.Context, tenantID int64) error
}

func (f *fakeListCache) GetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, bool, error) {
	if f.getConversationListFn != nil {
		return f.getConversationListFn(ctx, tenantID, filter)
	}
	return nil, false, nil
}

func (f *fakeListCache) SetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter, value []model.Conversation) error {
	if f.setConversationListFn != nil {
		return f.setConversationListFn(ctx, tenantID, filter, value)
	}
	return nil
}

func (f *fakeListCache) InvalidateConversationList(ctx context.Context, tenantID int64) error {
	if f.invalidateConversationFn != nil {
		return f.invalidateConversationFn(ctx, tenantID)
	}
	return nil
}

func (f *fakeListCache) GetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, bool, error) {
	if f.getTicketListFn != nil {
		return f.getTicketListFn(ctx, tenantID, filter)
	}
	return nil, false, nil
}

func (f *fakeListCache) SetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter, value []model.Ticket) error {
	if f.setTicketListFn != nil {
		return f.setTicketListFn(ctx, tenantID, filter, value)
	}
	return nil
}

func (f *fakeListCache) InvalidateTicketList(ctx context.Context, tenantID int64) error {
	if f.invalidateTicketFn != nil {
		return f.invalidateTicketFn(ctx, tenantID)
	}
	return nil
}

type recordingDispatcher struct {
	events []model.DomainEvent
	err    error
}

func (d *recordingDispatcher) Dispatch(ctx context.Context, event model.DomainEvent) error {
	if d.err != nil {
		return d.err
	}
	d.events = append(d.events, event)
	return nil
}
