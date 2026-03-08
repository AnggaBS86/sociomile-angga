package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"sociomile-be/internal/domain/model"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) FindActiveByCustomerID(ctx context.Context, tenantID, customerID int64) (*model.Conversation, error) {
	q := `
		SELECT id, tenant_id, customer_id, status, assigned_agent_id, created_at, updated_at
		FROM conversations
		WHERE tenant_id = ? AND customer_id = ? AND status IN ('open', 'assigned')
		ORDER BY id DESC
		LIMIT 1`

	var c model.Conversation
	var assigned sql.NullInt64
	err := r.db.QueryRowContext(ctx, q, tenantID, customerID).Scan(&c.ID, &c.TenantID, &c.CustomerID, &c.Status, &assigned, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if assigned.Valid {
		c.AssignedAgentID = &assigned.Int64
	}

	return &c, nil
}

func (r *ConversationRepository) Create(ctx context.Context, conversation *model.Conversation) error {
	res, err := r.db.ExecContext(ctx, `INSERT INTO conversations (tenant_id, customer_id, status, assigned_agent_id) VALUES (?, ?, ?, ?)`,
		conversation.TenantID,
		conversation.CustomerID,
		conversation.Status,
		conversation.AssignedAgentID,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	conversation.ID = id

	return nil
}

func (r *ConversationRepository) AddMessage(ctx context.Context, message *model.Message) error {
	res, err := r.db.ExecContext(ctx, `INSERT INTO messages (conversation_id, sender_type, message) VALUES (?, ?, ?)`,
		message.ConversationID,
		message.SenderType,
		message.Message,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	message.ID = id

	return nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, tenantID, conversationID int64) (*model.Conversation, error) {
	q := `SELECT id, tenant_id, customer_id, status, assigned_agent_id, created_at, updated_at
		FROM conversations
		WHERE tenant_id = ? AND id = ? LIMIT 1`
	var c model.Conversation
	var assigned sql.NullInt64
	err := r.db.QueryRowContext(ctx, q, tenantID, conversationID).Scan(&c.ID, &c.TenantID, &c.CustomerID, &c.Status, &assigned, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if assigned.Valid {
		c.AssignedAgentID = &assigned.Int64
	}

	return &c, nil
}

func (r *ConversationRepository) List(ctx context.Context, tenantID int64, filter model.ConversationFilter) ([]model.Conversation, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	parts := []string{"SELECT id, tenant_id, customer_id, status, assigned_agent_id, created_at, updated_at FROM conversations WHERE tenant_id = ?"}
	args := []any{tenantID}

	if filter.Status != "" {
		parts = append(parts, "AND status = ?")
		args = append(args, filter.Status)
	}

	if filter.AssignedAgentID != nil {
		parts = append(parts, "AND assigned_agent_id = ?")
		args = append(args, *filter.AssignedAgentID)
	}

	parts = append(parts, "ORDER BY id DESC LIMIT ? OFFSET ?")
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, strings.Join(parts, " "), args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	conversations := make([]model.Conversation, 0)
	for rows.Next() {
		var c model.Conversation
		var assigned sql.NullInt64
		if err := rows.Scan(&c.ID, &c.TenantID, &c.CustomerID, &c.Status, &assigned, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}

		if assigned.Valid {
			c.AssignedAgentID = &assigned.Int64
		}

		conversations = append(conversations, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

func (r *ConversationRepository) ListMessages(ctx context.Context, conversationID int64) ([]model.Message, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, conversation_id, sender_type, message, created_at FROM messages WHERE conversation_id = ? ORDER BY id ASC`, conversationID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := make([]model.Message, 0)
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderType, &m.Message, &m.CreatedAt); err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *ConversationRepository) UpdateAssignment(ctx context.Context, conversationID int64, assignedAgentID int64, status string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE conversations SET assigned_agent_id = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, assignedAgentID, status, conversationID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("conversation not updated")
	}

	return nil
}
