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
	if contCheckRes, message, containerID := s.docker.ContainerExist(ctx, userID); contCheckRes {
		fmt.Println(message)
		s.docker.Stop(ctx, containerID)
		s.docker.Remove(ctx, containerID)
	}

	containerID := s.docker.Create(ctx, userID)
	s.docker.Start(ctx, containerID)
	s.cache.Set(containerID, time.Now().GoString())

	go func() {
		fmt.Println("container stopped")
	}()

	s.docker.Attach(ctx, containerID)
}

func (s *Session) CleanupContainer() {
	// ctx := context.Background()
	fmt.Println("stopping containers")
	// s.docker.Stop(ctx, containerID)
	// s.docker.Remove(ctx, containerID)

}
