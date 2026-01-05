// Package transport
package transport

import (
	"context"
	"io"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *SessionHandler) Terminal(c echo.Context) error {
	userID := c.QueryParam("userID")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "userID query parameter is required",
		})
	}

	ctx := context.Background()
	exists, _, containerID := h.svc.GetDocker().ContainerExist(ctx, userID)
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "container not found for userID: " + userID,
		})
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	reader, writer, err := h.svc.GetDocker().AttachToContainer(ctx, containerID)
	if err != nil {
		c.Logger().Errorf("Failed to attach to container: %v", err)
		return err
	}
	defer reader.Close()
	defer writer.Close()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err != io.EOF {
					c.Logger().Errorf("Error reading from container: %v", err)
				}
				break
			}
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
					c.Logger().Errorf("Error writing to websocket: %v", err)
					break
				}
			}
		}
	}()

	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger().Errorf("Websocket error: %v", err)
			}
			break
		}

		if messageType == websocket.BinaryMessage || messageType == websocket.TextMessage {
			if _, err := writer.Write(data); err != nil {
				c.Logger().Errorf("Error writing to container: %v", err)
				break
			}
		}
	}

	return nil
}
