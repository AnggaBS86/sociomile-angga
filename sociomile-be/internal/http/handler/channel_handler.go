package handler

import (
	"sociomile-be/internal/http/request"
	"sociomile-be/internal/http/response"
	"sociomile-be/internal/service"
	"sociomile-be/internal/service/ratelimiter"

	"github.com/labstack/echo/v4"
)

type ChannelHandler struct {
	conversationService *service.ConversationService
	rateLimiter         ratelimiter.WebhookRateLimiter
}

func NewChannelHandler(conversationService *service.ConversationService, rateLimiter ratelimiter.WebhookRateLimiter) *ChannelHandler {
	return &ChannelHandler{conversationService: conversationService, rateLimiter: rateLimiter}
}

func (h *ChannelHandler) Webhook(c echo.Context) error {
	var req request.ChannelWebhookRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("invalid payload")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationFailed(err)
	}

	allowed, err := h.rateLimiter.Allow(c.Request().Context(), req.TenantID)
	if err != nil {
		return response.InternalServerError("failed to enforce webhook rate limit")
	}
	if !allowed {
		return response.TooManyRequests("webhook rate limit exceeded")
	}

	conversation, err := h.conversationService.IngestChannelMessage(c.Request().Context(), req.TenantID, req.CustomerExternalID, req.Message)
	if err != nil {
		if err == service.ErrInvalidInput {
			return response.BadRequest(err.Error())
		}

		return response.InternalServerError("failed to ingest message")
	}

	payload := map[string]any{
		"conversation_id": conversation.ID,
		"status":          conversation.Status,
	}

	return response.OK(c, "operation succeeded", payload)
}
