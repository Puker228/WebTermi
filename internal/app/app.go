// Package app solve project in one func
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	transport.MiddlewareRegister(e)
	transport.RouterRegister(e, handler)

	e.Static("/", "./public")

	// gracefull shutdown. Ref: https://echo.labstack.com/docs/cookbook/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
