package session

import (
	"context"
	"fmt"

	"github.com/Puker228/WebTermi/internal/docker"
)

func StartSession(dockerService docker.Service, userID string) {
	fmt.Println("starting session")
	fmt.Println("create and start container")
	ctx := context.Background()
	containerID := dockerService.CreateAndStart(ctx, "user-"+userID)
	dockerService.Attach(ctx, containerID)
	fmt.Println("stopping container")
	dockerService.RemoveContainer(containerID)
}
