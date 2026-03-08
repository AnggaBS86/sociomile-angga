package events

import (
	"context"
	"encoding/json"

	"sociomile-be/internal/domain/model"
	repository "sociomile-be/internal/domain/repository_interface"

	"github.com/hibiken/asynq"
)

const (
	EventConversationAssigned  = "conversation.assigned"
	EventConversationEscalated = "conversation.escalated"
	EventTicketCreated         = "ticket.created"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, event model.DomainEvent) error
}

type AsynqDispatcher struct {
	activityRepo repository.ActivityLogRepository
	client       *asynq.Client
	queue        string
}

func NewAsynqDispatcher(activityRepo repository.ActivityLogRepository, redisOpt asynq.RedisClientOpt, queue string) *AsynqDispatcher {
	if queue == "" {
		queue = "events"
	}

	return &AsynqDispatcher{
		activityRepo: activityRepo,
		client:       asynq.NewClient(redisOpt),
		queue:        queue,
	}
}

func (d *AsynqDispatcher) Dispatch(ctx context.Context, event model.DomainEvent) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	activity := &model.ActivityLog{
		TenantID:   event.TenantID,
		EntityType: event.EntityType,
		EntityID:   event.EntityID,
		EventName:  event.EventName,
		Payload:    string(payloadBytes),
	}
	if err := d.activityRepo.Create(ctx, activity); err != nil {
		return err
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	task := asynq.NewTask(event.EventName, eventBytes)
	_, err = d.client.Enqueue(task,
		asynq.Queue(d.queue),
		asynq.MaxRetry(10),
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *AsynqDispatcher) Close() error {
	if d.client == nil {
		return nil
	}
	return d.client.Close()
}

type NoopDispatcher struct{}

func NewNoopDispatcher() *NoopDispatcher {
	return &NoopDispatcher{}
}

func (d *NoopDispatcher) Dispatch(ctx context.Context, event model.DomainEvent) error {
	return nil
}
