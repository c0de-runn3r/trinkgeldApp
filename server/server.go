package server

import (
	"trinkgeldApp/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeEchoServer(h *handlers.AppContext) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")

	// API requests

	api.POST("/submit-gastromatic-report", h.SubmitGastromatic)
	api.POST("/submit-gastronovi-report", h.SubmitGastronovi)
	api.GET("/get-tips-per-day", h.GetTipsPerDay)
	api.GET("/calculate-tips", h.CalculateTips)

	return e
}
