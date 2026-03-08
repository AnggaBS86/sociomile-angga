package asynqworker

import (
	"context"

	repository "sociomile-be/internal/domain/repository_interface"
	"sociomile-be/internal/service/cache"
	"sociomile-be/internal/service/events"

	"github.com/hibiken/asynq"
)

type Worker struct {
	server           *asynq.Server
	mux              *asynq.ServeMux
	conversationRepo repository.ConversationRepository
	ticketRepo       repository.TicketRepository
	listCache        cache.ListCache
	dispatcher       events.Dispatcher
}

func New(
	redisOpt asynq.RedisClientOpt,
	queue string,
	concurrency int,
	conversationRepo repository.ConversationRepository,
	ticketRepo repository.TicketRepository,
	listCache cache.ListCache,
	dispatcher events.Dispatcher,
) *Worker {
	if queue == "" {
		queue = "events"
	}
	if concurrency <= 0 {
		concurrency = 10
	}

	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: concurrency,
		Queues:      map[string]int{queue: 1},
	})

	w := &Worker{
		server:           server,
		mux:              asynq.NewServeMux(),
		conversationRepo: conversationRepo,
		ticketRepo:       ticketRepo,
		listCache:        listCache,
		dispatcher:       dispatcher,
	}

	w.mux.HandleFunc(events.EventConversationAssigned, w.handleConversationAssigned)
	w.mux.HandleFunc(events.EventConversationEscalated, w.handleConversationEscalated)
	w.mux.HandleFunc(events.EventTicketCreated, w.handleTicketCreated)

	return w
}

func (w *Worker) Run() error {
	return w.server.Run(w.mux)
}

func (w *Worker) Shutdown(ctx context.Context) {
	w.server.Shutdown()
}
