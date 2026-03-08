package model

import "time"

const (
	RoleAdmin = "admin"
	RoleAgent = "agent"
)

type User struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
