package main

import (
	"log"
	"time"

	"sociomile-be/internal/config"
	"sociomile-be/internal/http/handler"
	httpmiddleware "sociomile-be/internal/http/middleware"
	"sociomile-be/internal/http/router"
	"sociomile-be/internal/platform/database"
	redisplatform "sociomile-be/internal/platform/redis"
	repositories "sociomile-be/internal/repository"
	"sociomile-be/internal/service"
	"sociomile-be/internal/service/cache"
	"sociomile-be/internal/service/events"
	"sociomile-be/internal/service/ratelimiter"

	"github.com/go-playground/validator/v10"
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

	userRepo := repositories.NewUserRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	activityLogRepo := repositories.NewActivityLogRepository(db)

	listCache := cache.NewRedisListCache(
		redisClient,
		time.Duration(cfg.ConversationCacheTTL)*time.Second,
		time.Duration(cfg.TicketCacheTTL)*time.Second,
	)

	asynqRedisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	eventDispatcher := events.NewAsynqDispatcher(activityLogRepo, asynqRedisOpt, cfg.AsynqQueue)
	defer func() {
		if err := eventDispatcher.Close(); err != nil {
			log.Printf("failed to close asynq client: %v", err)
		}
	}()

	authService := service.NewAuthService(userRepo)
	conversationService := service.NewConversationService(customerRepo, conversationRepo, listCache, eventDispatcher)
	ticketService := service.NewTicketService(conversationRepo, ticketRepo, listCache, eventDispatcher)

	webhookRateLimiter := ratelimiter.NewRedisWebhookRateLimiter(redisClient, cfg.WebhookRateLimitPerMinute)

	authHandler := handler.NewAuthHandler(authService, cfg.JWTSecret, cfg.JWTTTLMinutes)
	channelHandler := handler.NewChannelHandler(conversationService, webhookRateLimiter)
	conversationHandler := handler.NewConversationHandler(conversationService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	authMW := httpmiddleware.NewAuthMiddleware(cfg.JWTSecret)

	e := router.New(authMW, authHandler, channelHandler, conversationHandler, ticketHandler)
	e.Validator = &httpmiddleware.CustomValidator{Validator: validator.New()}

	log.Printf("server listening on :%s", cfg.AppPort)
	if err := e.Start(":" + cfg.AppPort); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
