package handler

import (
	"strconv"

	"sociomile-be/internal/domain/model"
	httpmiddleware "sociomile-be/internal/http/middleware"
	"sociomile-be/internal/http/request"
	"sociomile-be/internal/http/response"
	"sociomile-be/internal/service"

	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) Escalate(c echo.Context) error {
	claims, ok := httpmiddleware.GetClaims(c)
	if !ok {
		return response.Unauthorized("missing auth")
	}

	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest("invalid conversation id")
	}

	var req request.EscalateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("invalid payload")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationFailed(err)
	}

	ticket, err := h.ticketService.Escalate(c.Request().Context(), claims.TenantID, conversationID, claims.UserID, req.Title, req.Description, req.Priority)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			return response.NotFound(err.Error())
		case service.ErrAlreadyExists:
			return response.Conflict("conversation already escalated")
		default:
			return response.InternalServerError("failed to escalate")
		}
	}

	return response.Accepted(c, "event queued", ticket)
}

func (h *TicketHandler) UpdateStatus(c echo.Context) error {
	claims, ok := httpmiddleware.GetClaims(c)
	if !ok {
		return response.Unauthorized("missing auth")
	}

	ticketID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest("invalid ticket id")
	}

	var req request.UpdateTicketStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("invalid payload")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationFailed(err)
	}

	if err := h.ticketService.UpdateStatus(c.Request().Context(), claims.TenantID, ticketID, req.Status); err != nil {
		if err == service.ErrInvalidInput {
			return response.BadRequest(err.Error())
		}
		return response.InternalServerError("failed to update ticket")
	}

	return response.OK(c, "operation succeeded", map[string]string{"status": "ok"})
}

func (h *TicketHandler) List(c echo.Context) error {
	claims, ok := httpmiddleware.GetClaims(c)
	if !ok {
		return response.Unauthorized("missing auth")
	}

	filter := model.TicketFilter{
		Status: c.QueryParam("status"),
		Limit:  parseInt(c.QueryParam("limit"), 20),
		Offset: parseInt(c.QueryParam("offset"), 0),
	}

	if v := c.QueryParam("assigned_agent"); v != "" {
		agentID := int64(parseInt(v, 0))
		if agentID > 0 {
			filter.AssignedAgentID = &agentID
		}
	}

	tickets, err := h.ticketService.List(c.Request().Context(), claims.TenantID, filter)
	if err != nil {
		return response.InternalServerError("failed to list tickets")
	}

	return response.OK(c, "operation succeeded", tickets)
}
