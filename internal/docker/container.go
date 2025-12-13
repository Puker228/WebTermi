package docker

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Service struct {
	client *client.Client
}

func NewService(client *client.Client) *Service {
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

func (s *Service) RemoveContainer(containerID string) {
	log.Fatalf("Removing container:%v", containerID)
}
