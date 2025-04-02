package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/app"
	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Загрузка переменных окружения
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	// Конфигурация
	cfg := config.NewConfig("./config/config.json")

	bg := context.Background()
	ctx, stop := signal.NotifyContext(bg, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := app.New(cfg)
	go app.Run(ctx)

	<-ctx.Done()

	fmt.Println("Shutdown")
	log.Println("Shutdown")

	shutdownCtx, cancel := context.WithTimeout(bg, time.Second*5)
	defer cancel()

	err = app.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %s", err)
	}
}
