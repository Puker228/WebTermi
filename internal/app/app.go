// Package app solve project in one func
package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Puker228/WebTermi/internal/cache"
	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
	"github.com/Puker228/WebTermi/internal/transport"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron"
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

	c := cron.New()
	c.AddFunc("@every 1m", func() {
		ctx := context.Background()
		res, err := dockerSvc.ContainerList(ctx)
		if err != nil {
			fmt.Println("228")
		}
		fmt.Println(res)
	})
	c.Start()

	e.Logger.Fatal(e.Start(":1323"))
}
