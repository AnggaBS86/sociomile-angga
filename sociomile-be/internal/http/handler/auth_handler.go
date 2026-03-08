package handler

import (
	httpmiddleware "sociomile-be/internal/http/middleware"
	"sociomile-be/internal/http/request"
	"sociomile-be/internal/http/response"
	"sociomile-be/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
	jwtTTL      int
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string, jwtTTL int) *AuthHandler {
	return &AuthHandler{authService: authService, jwtSecret: jwtSecret, jwtTTL: jwtTTL}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("invalid payload")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationFailed(err)
	}

	user, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidInput, service.ErrInvalidCredential:
			return response.Unauthorized(err.Error())
		default:
			return response.InternalServerError("failed to login")
		}
	}

	token, err := httpmiddleware.GenerateToken(h.jwtSecret, user.ID, user.TenantID, user.Role, h.jwtTTL)
	if err != nil {
		return response.InternalServerError("failed to create token")
	}

	payload := map[string]any{
		"token": token,
		"user": map[string]any{
			"id":        user.ID,
			"tenant_id": user.TenantID,
			"email":     user.Email,
			"role":      user.Role,
		},
	}

	return response.OK(c, "operation succeeded", payload)
}
