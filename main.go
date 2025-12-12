package main

import (
	"context"
	"io"
	"os"

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

	containerExist(containerName, ctx, apiClient)

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

func containerExist(containerName string, ctx context.Context, apiClient *client.Client) bool {
	cont_list, err := apiClient.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	cont_check := "/" + containerName

	for _, cont := range cont_list.Items {
		for _, name := range cont.Names {
			if cont_check == name {
				return true
			}
		}
	}

	return false
}
