// Package session
package session

import (
	"context"
	"fmt"
	"time"
)

type Session struct {
	docker ContainerRuntime
	cache  UserCache
}

func NewSessionService(dockerPort ContainerRuntime, cache UserCache) *Session {
	return &Session{docker: dockerPort, cache: cache}
}

func (s *Session) StartSession(userID string) {
	fmt.Println("starting session")
	fmt.Println("create and start container")
	ctx := context.Background()
	if contCheckRes := s.docker.ContainerExist(ctx, userID); contCheckRes.Exist {
		s.docker.Stop(ctx, contCheckRes.ContainerID)
		s.docker.Remove(ctx, contCheckRes.ContainerID)
	}

	containerID := s.docker.Create(ctx, userID)
	s.docker.Start(ctx, containerID)
	s.cache.Set(containerID, time.Now().GoString())

	go func() {
		fmt.Println("stopping container")
		<-time.After(20 * time.Second)
		s.docker.Stop(ctx, containerID)
		s.docker.Remove(ctx, containerID)
		fmt.Println("container stopped")
	}()

	s.docker.Attach(ctx, containerID)
}
