package model

import "time"

type Ticket struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	ConversationID  int64     `json:"conversation_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	Priority        string    `json:"priority"`
	AssignedAgentID *int64    `json:"assigned_agent_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TicketFilter struct {
	Status          string
	AssignedAgentID *int64
	Limit           int
	Offset          int
}
