// Package app
package app

import (
	"log"

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
	defer dockerClient.Close()

	redisClient := cache.NewClient()
	defer redisClient.Close()

	dockerSvc := docker.NewContainerService(dockerClient)
	redisSVC := cache.NewRedisService(redisClient)

	sessionService := session.NewSessionService(dockerSvc, redisSVC)

	handler := transport.NewSessionHandler(sessionService)

	e := echo.New()
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/session", handler.Start)

	e.Logger.Fatal(e.Start(":1323"))
}
