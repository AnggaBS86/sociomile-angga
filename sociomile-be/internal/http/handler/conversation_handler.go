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

type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler(conversationService *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{conversationService: conversationService}
}

func (h *ConversationHandler) List(c echo.Context) error {
	claims, ok := httpmiddleware.GetClaims(c)
	if !ok {
		return response.Unauthorized("missing auth")
	}

	filter := model.ConversationFilter{
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

	conversations, err := h.conversationService.List(c.Request().Context(), claims.TenantID, filter)
	if err != nil {
		return response.InternalServerError("failed to list conversations")
	}

	return response.OK(c, "operation succeeded", conversations)
}

func (h *ConversationHandler) Detail(c echo.Context) error {
	claims, ok := httpmiddleware.GetClaims(c)
	if !ok {
		return response.Unauthorized("missing auth")
	}

	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest("invalid conversation id")
	}

	detail, err := h.conversationService.GetDetail(c.Request().Context(), claims.TenantID, conversationID)
	if err != nil {
		if err == service.ErrNotFound {
			return response.NotFound("conversation not found")
		}
		return response.InternalServerError("failed to get conversation")
	}

	return response.OK(c, "operation succeeded", detail)
}

func (h *ConversationHandler) Reply(c echo.Context) error {
	claims, ok := httpmiddleware.GetClaims(c)
	if !ok {
		return response.Unauthorized("missing auth")
	}

	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest("invalid conversation id")
	}

	var req request.ReplyRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("invalid payload")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationFailed(err)
	}

	err = h.conversationService.AgentReply(c.Request().Context(), claims.TenantID, conversationID, claims.UserID, claims.Role, req.Message)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			return response.BadRequest(err.Error())
		case service.ErrForbidden:
			return response.Forbidden(err.Error())
		case service.ErrNotFound:
			return response.NotFound(err.Error())
		default:
			return response.InternalServerError("failed to send reply")
		}
	}

	return response.OK(c, "operation succeeded", map[string]string{"status": "ok"})
}

func parseInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return v
}
