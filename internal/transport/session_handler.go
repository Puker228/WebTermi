// Package transport
package transport

import (
	"net/http"

	"github.com/Puker228/WebTermi/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	svc *session.Session
}

func NewSessionHandler(svc *session.Session) *SessionHandler {
	return &SessionHandler{svc: svc}
}

func (h *SessionHandler) Start(c echo.Context) error {
	userID := uuid.NewString()

	go h.svc.StartSession(userID)

	status := "session " + userID + " started"

	return c.JSON(http.StatusOK, map[string]string{
		"status": status,
	})
}
