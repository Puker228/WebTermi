package session

import (
	"context"

	"github.com/Puker228/WebTermi/internal/docker"
)

type ContainerPort interface {
	CreateAndStart(ctx context.Context, containerName string) string
	Attach(ctx context.Context, containerID string)
	RemoveContainer(ctx context.Context, containerID string)
	ContainerExist(ctx context.Context, containerName string) docker.ContainerCheckResult
}
