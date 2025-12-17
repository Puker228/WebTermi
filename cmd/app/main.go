package main

import (
	"fmt"
	"log"

	"github.com/Puker228/WebTermi/internal/app"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	app.RunServer()
	fmt.Println("start cron")
}
