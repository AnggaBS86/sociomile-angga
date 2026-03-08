package model

import "time"

type ActivityLog struct {
	ID         int64      `json:"id"`
	TenantID   int64      `json:"tenant_id"`
	EntityType string     `json:"entity_type"`
	EntityID   int64      `json:"entity_id"`
	EventName  string     `json:"event_name"`
	Payload    string     `json:"payload"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

type DomainEvent struct {
	TenantID   int64          `json:"tenant_id"`
	EntityType string         `json:"entity_type"`
	EntityID   int64          `json:"entity_id"`
	EventName  string         `json:"event_name"`
	Payload    map[string]any `json:"payload"`
}
