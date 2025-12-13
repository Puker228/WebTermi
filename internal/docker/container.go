package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Service struct {
	client *client.Client
}

func NewContainerService(client *client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) CreateAndStart(ctx context.Context, containerName string) string {
	// pull image
	reader, err := s.client.ImagePull(ctx, "ubuntu:24.04", client.ImagePullOptions{})
	if err != nil {
		log.Fatalf("Failed to pull image: %v", err)
	}
	io.Copy(os.Stdout, reader)

	// create container
	resp, err := s.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:  containerName,
		Image: "ubuntu:24.04",
		Config: &container.Config{
			Cmd:          []string{"/bin/bash"},
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			OpenStdin:    true,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// starting container
	if _, err := s.client.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}
	return resp.ID
}

func (s *Service) Attach(ctx context.Context, containerID string) {
	conn, err := s.client.ContainerAttach(ctx, containerID, client.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		panic(err)
	}

	// output
	go io.Copy(os.Stdout, conn.Reader)

	// input
	go io.Copy(conn.Conn, os.Stdin)
}

func (s *Service) RemoveContainer(ctx context.Context, containerID string) {
	result, err := s.client.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		log.Fatalf("Failed to remove container: %v", err)
	}
	fmt.Println(result)
}

func (s *Service) ContainerExist(ctx context.Context, containerName string) ContainerCheckResult {
	contList, err := s.client.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	contCheck := "/" + containerName
	for _, cont := range contList.Items {
		if slices.Contains(cont.Names, contCheck) {
			return ContainerCheckResult{Exist: true, Message: ContainerExists, ContainerID: cont.ID}
		}
	}
	return ContainerCheckResult{Exist: false, Message: ContainerNotFound, ContainerID: ""}
}
