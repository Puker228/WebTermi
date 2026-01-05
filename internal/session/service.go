// Package session
package session

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	docker ContainerRuntime
	cache  UserCache
}

func NewSessionService(dockerPort ContainerRuntime, cache UserCache) *Service {
	return &Service{docker: dockerPort, cache: cache}
}

func (s *Service) GetDocker() ContainerRuntime {
	return s.docker
}

func (s *Service) StartSession(userID string) {
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
	if err := s.cache.Set(containerID, time.Now().GoString()); err != nil {
		fmt.Println("error while set key in redis")
	}

	fmt.Printf("Container %s created and started for userID %s\n", containerID, userID)
}
