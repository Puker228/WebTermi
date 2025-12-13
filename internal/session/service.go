package session

import (
	"context"
	"fmt"
	"time"

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
	if contCheckRes := s.docker.ContainerExist(ctx, userID); contCheckRes.Exist {
		s.docker.RemoveContainer(ctx, userID)
	}
	containerID := s.docker.CreateAndStart(ctx, userID)
	s.docker.Attach(ctx, containerID)

	func() {
		fmt.Println("stopping container")
		<-time.After(10 * time.Second)
		s.docker.RemoveContainer(ctx, containerID)
	}()
}
