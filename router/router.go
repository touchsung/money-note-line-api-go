package router

import (
	"github.com/Touchsung/money-note-line-api-go/handler"
	"github.com/labstack/echo/v4"
)

// InitRoute use for init all route in service
func Router() *echo.Echo{
	api := echo.New()


	// Middleware
	// api.Use(middleware.Logger())
	// api.Use(middleware.Recover())

	// Routes
  	api.GET("/", handler.Hello)

	// Line API
  	api.POST("/callback", handler.LineCallbackHandler)

	return api
}