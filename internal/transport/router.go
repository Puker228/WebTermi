package transport

import (
	"github.com/labstack/echo/v4"
)

func RouterRegister(e *echo.Echo, handler *SessionHandler) {
	// groupping api
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/session", handler.Start)
	apiV1.GET("/ws", handler.Hello)
}
