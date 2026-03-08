package main

import (
	"log"
	"time"

	"sociomile-be/internal/config"
	"sociomile-be/internal/platform/database"
	redisplatform "sociomile-be/internal/platform/redis"
	repositories "sociomile-be/internal/repository"
	"sociomile-be/internal/service/cache"
	"sociomile-be/internal/service/events"
	asynqworker "sociomile-be/internal/worker/asynq"

	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMySQL(cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	redisClient, err := redisplatform.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}
	}()

	conversationRepo := repositories.NewConversationRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	activityLogRepo := repositories.NewActivityLogRepository(db)

	listCache := cache.NewRedisListCache(
		redisClient,
		time.Duration(cfg.ConversationCacheTTL)*time.Second,
		time.Duration(cfg.TicketCacheTTL)*time.Second,
	)

	redisOpt := asynq.RedisClientOpt{Addr: cfg.RedisAddr, Password: cfg.RedisPassword, DB: cfg.RedisDB}
	dispatcher := events.NewAsynqDispatcher(activityLogRepo, redisOpt, cfg.AsynqQueue)
	defer func() {
		if err := dispatcher.Close(); err != nil {
			log.Printf("failed to close asynq dispatcher: %v", err)
		}
	}()

	worker := asynqworker.New(redisOpt, cfg.AsynqQueue, cfg.AsynqConcurrency, conversationRepo, ticketRepo, listCache, dispatcher)
	log.Printf("asynq worker started: queue=%s concurrency=%d", cfg.AsynqQueue, cfg.AsynqConcurrency)
	if err := worker.Run(); err != nil {
		log.Fatalf("asynq worker stopped: %v", err)
	}
}
