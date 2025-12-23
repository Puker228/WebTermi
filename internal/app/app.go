// Package app solve project in one func
package app

import (
	"log"

	"github.com/Puker228/WebTermi/internal/cache"
	"github.com/Puker228/WebTermi/internal/cron"
	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
	"github.com/Puker228/WebTermi/internal/transport"

	"github.com/labstack/echo/v4"
)

func RunServer() {
	// init clients
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	redisClient := cache.NewClient()
	defer dockerClient.Close()
	defer redisClient.Close()

	// init services and cron
	dockerSvc := docker.NewContainerService(dockerClient)
	redisSVC := cache.NewRedisService(redisClient)
	sessionService := session.NewSessionService(dockerSvc, redisSVC)

	cron.CleanUpCrone(dockerSvc)

	// init web server
	e := echo.New()

	handler := transport.NewSessionHandler(sessionService)
	transport.RouterRegister(e, handler)
	transport.MiddlewareRegister(e)

	e.Static("/", "./public")
	e.Logger.Fatal(e.Start(":1323"))
}
