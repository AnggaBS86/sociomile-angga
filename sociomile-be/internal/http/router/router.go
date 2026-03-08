package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"sociomile-be/internal/domain/model"
	"sociomile-be/internal/http/handler"
	httpmiddleware "sociomile-be/internal/http/middleware"
	"sociomile-be/internal/http/response"
)

func New(
	authMW *httpmiddleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
	channelHandler *handler.ChannelHandler,
	conversationHandler *handler.ConversationHandler,
	ticketHandler *handler.TicketHandler,
) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = httpmiddleware.HTTPErrorHandler

	e.GET("/health", func(c echo.Context) error {
		return response.OK(c, "operation succeeded", map[string]string{"status": "ok"})
	})

	e.POST("/auth/login", authHandler.Login)
	e.POST("/channel/webhook", channelHandler.Webhook)

	api := e.Group("", authMW.Authenticate)
	api.GET("/conversations", conversationHandler.List, httpmiddleware.RequireRoles(model.RoleAdmin, model.RoleAgent))
	api.GET("/conversations/:id", conversationHandler.Detail, httpmiddleware.RequireRoles(model.RoleAdmin, model.RoleAgent))
	api.POST("/conversations/:id/messages", conversationHandler.Reply, httpmiddleware.RequireRoles(model.RoleAdmin, model.RoleAgent))

	api.POST("/conversations/:id/escalate", ticketHandler.Escalate, httpmiddleware.RequireRoles(model.RoleAgent))
	api.GET("/tickets", ticketHandler.List, httpmiddleware.RequireRoles(model.RoleAdmin, model.RoleAgent))
	api.PATCH("/tickets/:id/status", ticketHandler.UpdateStatus, httpmiddleware.RequireRoles(model.RoleAdmin))

	return e
}
