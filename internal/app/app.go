// Package app solve project in one func
package app

import (
	"context"
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
	c.AddFunc("@every 20m", func() {
		ctx := context.Background()
		names, err := dockerSvc.ContainerList(ctx)
		if err != nil {
			log.Printf("cron: failed to list containers: %v", err)
			return
		}

		for _, name := range names {
			if name == "backend" || name == "redis" {
				continue
			}
			func(n string) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("cron: panic while stopping container %s: %v", n, r)
					}
				}()
				dockerSvc.Stop(ctx, n)
				log.Printf("cron: stopped container %s", n)
			}(name)
		}
	})
	c.Start()

	e.Logger.Fatal(e.Start(":1323"))
}
