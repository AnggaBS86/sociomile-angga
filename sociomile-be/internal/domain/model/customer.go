package model

import "time"

type Customer struct {
	ID         int64     `json:"id"`
	TenantID   int64     `json:"tenant_id"`
	ExternalID string    `json:"external_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
