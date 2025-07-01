package main

import (
	"auth-service/internal/api"
	"auth-service/internal/config"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/internal/token"
	"auth-service/internal/webhook"
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal(errors.Wrap(err, "Error loading .env file"))
		return
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Error processing config"))
		return
	}
	repo, err := repository.NewRepository(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error creating repository"))
		return
	}

	refreshManager := token.NewRefreshTokenManager(cfg.Token.LenToken)
	accessManager := token.NewJWTTokenManager([]byte(cfg.Token.Secret))
	web := webhook.NewWebhook("example.com")

	serviceInstance := service.NewService(repo, refreshManager, accessManager, web)

	app := api.NewRouters(&api.Routers{Service: serviceInstance}, accessManager)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		log.Printf("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutting down gracefully...")
}
