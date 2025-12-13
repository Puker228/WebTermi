package main

import (
	"context"
	"io"
	"log"
	"os"
	"slices"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func main() {
	deployContainer("test-1")
}

func deployContainer(containerName string) {
	ctx := context.Background()
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	if contCheckRes := containerExist(containerName, ctx, apiClient); contCheckRes.Exist {
		removeContainer(contCheckRes.ContainerID)
	}

	// подтягивание образа
	reader, err := apiClient.ImagePull(ctx, "ubuntu:24.04", client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	// создание контейнера
	resp, err := apiClient.ContainerCreate(ctx, client.ContainerCreateOptions{
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
		panic(err)
	}

	// запуск контейнера
	if _, err := apiClient.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// привязка терминала к контейнеру
	conn, err := apiClient.ContainerAttach(ctx, resp.ID, client.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		panic(err)
	}

	// вывод резульатата
	go io.Copy(os.Stdout, conn.Reader)

	// ввод
	io.Copy(conn.Conn, os.Stdin)
}

func containerExist(containerName string, ctx context.Context, apiClient *client.Client) ContainerCheckResult {
	contList, err := apiClient.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	cont_check := "/" + containerName

	for _, cont := range contList.Items {
		if slices.Contains(cont.Names, cont_check) {
			return ContainerCheckResult{Exist: true, Message: ContainerExists, ContainerID: cont.ID}
		}
	}
	return ContainerCheckResult{Exist: false, Message: ContainerNotFound, ContainerID: ""}
}

type ContainerCheckResult struct {
	Exist       bool
	Message     string
	ContainerID string
}

const (
	ContainerExists   = "Container exists"
	ContainerNotFound = "Container not found"
)

func removeContainer(contID string) {
	log.Fatalf("Container %v is exist", contID)
}
