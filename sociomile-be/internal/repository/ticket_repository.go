package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"sociomile-be/internal/domain/model"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) GetByConversationID(ctx context.Context, tenantID, conversationID int64) (*model.Ticket, error) {
	q := `SELECT id, tenant_id, conversation_id, title, description, status, priority, assigned_agent_id, created_at, updated_at FROM tickets WHERE tenant_id = ? AND conversation_id = ? LIMIT 1`
	var t model.Ticket
	var assigned sql.NullInt64
	err := r.db.QueryRowContext(ctx, q, tenantID, conversationID).Scan(
		&t.ID,
		&t.TenantID,
		&t.ConversationID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&assigned,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if assigned.Valid {
		t.AssignedAgentID = &assigned.Int64
	}

	return &t, nil
}

func (r *TicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	res, err := r.db.ExecContext(ctx, `INSERT INTO tickets (tenant_id, conversation_id, title, description, status, priority, assigned_agent_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ticket.TenantID,
		ticket.ConversationID,
		ticket.Title,
		ticket.Description,
		ticket.Status,
		ticket.Priority,
		ticket.AssignedAgentID,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ticket.ID = id

	return nil
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, tenantID, ticketID int64, status string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tickets SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE tenant_id = ? AND id = ?`, status, tenantID, ticketID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}

func (r *TicketRepository) List(ctx context.Context, tenantID int64, filter model.TicketFilter) ([]model.Ticket, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	parts := []string{"SELECT id, tenant_id, conversation_id, title, description, status, priority, assigned_agent_id, created_at, updated_at FROM tickets WHERE tenant_id = ?"}
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

	tickets := make([]model.Ticket, 0)
	for rows.Next() {
		var t model.Ticket
		var assigned sql.NullInt64

		if err := rows.Scan(&t.ID, &t.TenantID, &t.ConversationID, &t.Title, &t.Description, &t.Status, &t.Priority, &assigned, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if assigned.Valid {
			t.AssignedAgentID = &assigned.Int64
		}

		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}
