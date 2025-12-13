package session

import (
	"context"
	"fmt"

	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/moby/moby/client"
)

type Session struct {
	docker *docker.Service
}

func NewSessionService(client *client.Client) *Session {
	return &Session{
		docker: docker.NewContainerService(client),
	}
}

func (s *Session) StartSession(userID string) {
	fmt.Println("starting session")
	fmt.Println("create and start container")
	ctx := context.Background()
	containerID := s.docker.CreateAndStart(ctx, userID)
	s.docker.Attach(ctx, containerID)
	fmt.Println("stopping container")
	s.docker.RemoveContainer(containerID)
}
