// Package transport
package transport

import (
	"fmt"
	"net/http"

	"github.com/Puker228/WebTermi/internal/session"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	svc *session.Service
}

func NewSessionHandler(svc *session.Service) *SessionHandler {
	return &SessionHandler{svc: svc}
}

func (h *SessionHandler) Start(c echo.Context) error {
	userID := uuid.NewString()

	go h.svc.StartSession(userID)

	status := "session " + userID + " started"

	return c.JSON(http.StatusOK, map[string]string{
		"status": status,
		"userID": userID,
	})
}

var upgrader = websocket.Upgrader{}

func (h *SessionHandler) Hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
