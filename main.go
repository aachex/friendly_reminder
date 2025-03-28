package main

import (
	"log"

	"github.com/artemwebber1/friendly_reminder/internal/app"
	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/joho/godotenv"
)

func init() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
}

func main() {
	// Конфигурация
	cfg := config.NewConfig(`config\config.json`)

	app := app.New(cfg)
	app.Run()
}
