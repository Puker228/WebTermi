package app

import (
	"log"
	"net/http"

	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func Run() {
	apiClient, err := docker.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer apiClient.Close()
	dockerSvc := docker.NewContainerService(apiClient)

	sessionService := session.NewSessionService(dockerSvc)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		go sessionService.StartSession(uuid.NewString())
		return c.String(http.StatusOK, "Service Started")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
