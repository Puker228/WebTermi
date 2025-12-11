package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func main() {
	for i := range 10 {
		deployContainer(fmt.Sprintf("test-%v", i))
	}
}

func deployContainer(containerName string) {
	ctx := context.Background()
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

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

	go io.Copy(os.Stdout, conn.Reader)

	// TODO убрать go после тестов
	go io.Copy(conn.Conn, os.Stdin)
}
