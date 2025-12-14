package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Puker228/WebTermi/internal/cache"
	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/Puker228/WebTermi/internal/session"
	"github.com/google/uuid"
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
	fmt.Println(redisSVC)

	sessionService := session.NewSessionService(dockerSvc)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		go sessionService.StartSession(uuid.NewString())
		return c.String(http.StatusOK, "Service Started")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
