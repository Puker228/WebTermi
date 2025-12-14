package app

import (
	"log"

	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
	"github.com/google/uuid"
)

func Run() {
	apiClient, err := docker.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer apiClient.Close()
	dockerSvc := docker.NewContainerService(apiClient)

	sessionService := session.NewSessionService(dockerSvc)
	sessionService.StartSession(uuid.NewString())
}
