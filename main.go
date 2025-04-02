package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/app"
	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/pkg/graceful"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	// Конфигурация
	cfg := config.NewConfig(`config\config.json`)

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := app.New(cfg)
	go app.Run(ctx)

	graceful.WaitShutdown()

	fmt.Println("Shutdown")
	log.Println("Shutdown")

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err = app.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %s", err)
	}
}
