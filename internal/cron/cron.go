// Package cron
package cron

import (
	"context"
	"log"

	"github.com/Puker228/WebTermi/internal/docker"
	"github.com/robfig/cron"
)

func CleanUpCrone(dockerSvc *docker.Service) {
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
}
