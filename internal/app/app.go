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

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron"
)

var upgrader = websocket.Upgrader{}

func hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

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
	e.Static("/", "./public")
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":1323"))
}
