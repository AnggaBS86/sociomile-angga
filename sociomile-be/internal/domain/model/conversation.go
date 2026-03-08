package model

import "time"

type Conversation struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	CustomerID      int64     `json:"customer_id"`
	Status          string    `json:"status" validate:"required,oneof=open assigned closed"`
	AssignedAgentID *int64    `json:"assigned_agent_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Messages        []Message `json:"messages,omitempty"`
	Customer        *Customer `json:"customer,omitempty"`
	Ticket          *Ticket   `json:"ticket,omitempty"`
}

type Message struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	SenderType     string    `json:"sender_type"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

type ConversationFilter struct {
	Status          string
	AssignedAgentID *int64
	Limit           int
	Offset          int
}
