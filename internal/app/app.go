package app

import "github.com/Puker228/WebTermi/internal/docker"

func Run() {
	apiClient, err := docker.NewClient()
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()
}
