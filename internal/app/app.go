// Package app solve project in one func
package app

import (
	"fmt"
	"log"
	"time"

	"github.com/Puker228/WebTermi/internal/cache"
	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
	"github.com/Puker228/WebTermi/internal/transport"
	"github.com/labstack/echo/v4"
)

func RunServer() {
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	redisClient := cache.NewClient()
	defer dockerClient.Close()
	defer redisClient.Close()

	dockerSvc := docker.NewContainerService(dockerClient)
	redisSVC := cache.NewRedisService(redisClient)
	sessionService := session.NewSessionService(dockerSvc, redisSVC)

	handler := transport.NewSessionHandler(sessionService)

	e := echo.New()

	transport.RouterRegister(e, handler)

	e.Logger.Fatal(e.Start(":1323"))
}

func RunWorker() {
	for i := range 100 {
		if i%2 == 0 {
			time.Sleep(1 * time.Minute)
		}
		fmt.Printf("%v: WebTermiorker in running", i)
	}
}
