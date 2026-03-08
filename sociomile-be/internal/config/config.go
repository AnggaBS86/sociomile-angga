package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppPort                   string
	DBDSN                     string
	JWTSecret                 string
	JWTTTLMinutes             int
	RedisAddr                 string
	RedisPassword             string
	RedisDB                   int
	WebhookRateLimitPerMinute int
	ConversationCacheTTL      int
	TicketCacheTTL            int
	AsynqQueue                string
	AsynqConcurrency          int
}

func Load() Config {
	_ = loadDotEnv(".env")

	return Config{
		AppPort:                   getEnv("APP_PORT", "8080"),
		DBDSN:                     getEnv("DB_DSN", "root:root@tcp(localhost:3306)/sociomile?parseTime=true"),
		JWTSecret:                 getEnv("JWT_SECRET", "change-me"),
		JWTTTLMinutes:             getEnvInt("JWT_TTL_MINUTES", 120),
		RedisAddr:                 getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:             getEnv("REDIS_PASSWORD", ""),
		RedisDB:                   getEnvInt("REDIS_DB", 0),
		WebhookRateLimitPerMinute: getEnvInt("WEBHOOK_RATE_LIMIT_PER_MINUTE", 60),
		ConversationCacheTTL:      getEnvInt("CONVERSATION_CACHE_TTL_SECONDS", 30),
		TicketCacheTTL:            getEnvInt("TICKET_CACHE_TTL_SECONDS", 30),
		AsynqQueue:                getEnv("ASYNQ_QUEUE", "events"),
		AsynqConcurrency:          getEnvInt("ASYNQ_CONCURRENCY", 10),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var out int
	if _, err := fmt.Sscanf(v, "%d", &out); err != nil {
		return fallback
	}
	return out
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)

		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, val)
	}

	return s.Err()
}
