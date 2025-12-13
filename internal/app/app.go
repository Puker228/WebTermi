package app

import (
	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
)

func Run() {
	apiClient, err := docker.NewClient()
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	dockerSvc := docker.NewContainerService(apiClient)

	sessionService := session.NewSessionService(dockerSvc)
	sessionService.StartSession("user123")
}
