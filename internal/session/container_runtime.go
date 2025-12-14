package session

import (
	"context"

	"github.com/Puker228/WebTermi/internal/docker"
)

type ContainerRuntime interface {
	Create(ctx context.Context, containerName string) string
	Start(ctx context.Context, containerName string)
	Attach(ctx context.Context, containerID string)
	Remove(ctx context.Context, containerID string)
	Stop(ctx context.Context, containerID string)
	ContainerExist(ctx context.Context, containerName string) docker.ContainerCheckResult
}
