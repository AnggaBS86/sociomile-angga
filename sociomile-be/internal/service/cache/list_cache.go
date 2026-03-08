package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sociomile-be/internal/domain/model"

	goredis "github.com/redis/go-redis/v9"
)

type ListCache interface {
	GetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, bool, error)
	SetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter, value []model.Conversation) error
	InvalidateConversationList(ctx context.Context, tenantID int64) error
	GetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, bool, error)
	SetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter, value []model.Ticket) error
	InvalidateTicketList(ctx context.Context, tenantID int64) error
}

type RedisListCache struct {
	client          *goredis.Client
	conversationTTL time.Duration
	ticketTTL       time.Duration
}

func NewRedisListCache(client *goredis.Client, conversationTTL, ticketTTL time.Duration) *RedisListCache {
	return &RedisListCache{
		client:          client,
		conversationTTL: conversationTTL,
		ticketTTL:       ticketTTL,
	}
}

func (c *RedisListCache) GetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, bool, error) {
	if c == nil || c.client == nil {
		return nil, false, nil
	}
	raw, err := c.client.Get(ctx, conversationListKey(tenantID, filter)).Bytes()
	if err == goredis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var out []model.Conversation
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, false, err
	}
	return out, true, nil
}

func (c *RedisListCache) SetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter, value []model.Conversation) error {
	if c == nil || c.client == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, conversationListKey(tenantID, filter), raw, c.conversationTTL).Err()
}

func (c *RedisListCache) InvalidateConversationList(ctx context.Context, tenantID int64) error {
	return c.deleteByPattern(ctx, fmt.Sprintf("cache:conversations:list:tenant:%d:*", tenantID))
}

func (c *RedisListCache) GetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, bool, error) {
	if c == nil || c.client == nil {
		return nil, false, nil
	}
	raw, err := c.client.Get(ctx, ticketListKey(tenantID, filter)).Bytes()
	if err == goredis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var out []model.Ticket
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, false, err
	}
	return out, true, nil
}

func (c *RedisListCache) SetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter, value []model.Ticket) error {
	if c == nil || c.client == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, ticketListKey(tenantID, filter), raw, c.ticketTTL).Err()
}

func (c *RedisListCache) InvalidateTicketList(ctx context.Context, tenantID int64) error {
	return c.deleteByPattern(ctx, fmt.Sprintf("cache:tickets:list:tenant:%d:*", tenantID))
}

func (c *RedisListCache) deleteByPattern(ctx context.Context, pattern string) error {
	if c == nil || c.client == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func conversationListKey(tenantID int64, filter model.ConversationFilter) string {
	assigned := "none"
	if filter.AssignedAgentID != nil {
		assigned = fmt.Sprintf("%d", *filter.AssignedAgentID)
	}
	return fmt.Sprintf("cache:conversations:list:tenant:%d:status:%s:assigned:%s:limit:%d:offset:%d",
		tenantID, filter.Status, assigned, filter.Limit, filter.Offset)
}

func ticketListKey(tenantID int64, filter model.TicketFilter) string {
	assigned := "none"
	if filter.AssignedAgentID != nil {
		assigned = fmt.Sprintf("%d", *filter.AssignedAgentID)
	}
	return fmt.Sprintf("cache:tickets:list:tenant:%d:status:%s:assigned:%s:limit:%d:offset:%d",
		tenantID, filter.Status, assigned, filter.Limit, filter.Offset)
}

type NoopListCache struct{}

func NewNoopListCache() *NoopListCache { return &NoopListCache{} }

func (c *NoopListCache) GetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, bool, error) {
	return nil, false, nil
}
func (c *NoopListCache) SetConversationList(ctx context.Context, tenantID int64, filter model.ConversationFilter, value []model.Conversation) error {
	return nil
}
func (c *NoopListCache) InvalidateConversationList(ctx context.Context, tenantID int64) error { return nil }
func (c *NoopListCache) GetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, bool, error) {
	return nil, false, nil
}
func (c *NoopListCache) SetTicketList(ctx context.Context, tenantID int64, filter model.TicketFilter, value []model.Ticket) error {
	return nil
}
func (c *NoopListCache) InvalidateTicketList(ctx context.Context, tenantID int64) error { return nil }
