package session

import (
	"context"
	"fmt"
	"time"
)

type Session struct {
	docker ContainerPort
}

func NewSessionService(dockerPort ContainerPort) *Session {
	return &Session{docker: dockerPort}
}

func (s *Session) StartSession(userID string) {
	fmt.Println("starting session")
	fmt.Println("create and start container")
	ctx := context.Background()
	if contCheckRes := s.docker.ContainerExist(ctx, userID); contCheckRes.Exist {
		s.docker.RemoveContainer(ctx, contCheckRes.ContainerID)
	}
	containerID := s.docker.CreateAndStart(ctx, userID)

	go func() {
		fmt.Println("stopping container")
		<-time.After(20 * time.Second)
		s.docker.RemoveContainer(ctx, containerID)
	}()

	s.docker.Attach(ctx, containerID)
}
