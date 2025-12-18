package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Service struct {
	client *client.Client
}

func NewContainerService(client *client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, containerName string) string {
	// TODO добавить возможность выбрать образ
	// create container
	resp, err := s.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:  containerName,
		Image: "myub:latest",
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

	return resp.ID
}

func (s *Service) Start(ctx context.Context, containerID string) {
	if _, err := s.client.ContainerStart(ctx, containerID, client.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}
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
	io.Copy(conn.Conn, os.Stdin)
}

func (s *Service) Stop(ctx context.Context, containerID string) {
	if _, err := s.client.ContainerStop(ctx, containerID, client.ContainerStopOptions{}); err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}
}

func (s *Service) Remove(ctx context.Context, containerID string) {
	// remove container
	result, err := s.client.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		log.Fatalf("Failed to remove container: %v", err)
	}
	fmt.Println(result)
}

func (s *Service) ContainerExist(ctx context.Context, containerName string) (bool, string, string) {
	contList, err := s.client.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	contCheck := "/" + containerName
	for _, cont := range contList.Items {
		if slices.Contains(cont.Names, contCheck) {
			return true, ContainerExists, cont.ID
		}
	}
	return false, ContainerNotFound, ""
}

func (s *Service) ContainerList(ctx context.Context) ([]string, error) {
	contList, err := s.client.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	var result []string
	for _, cont := range contList.Items {
		for _, name := range cont.Names {
			result = append(result, strings.TrimPrefix(name, "/"))
		}
	}
	return result, nil
}
