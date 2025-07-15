package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oyavri/aldim_verdim/internal/frontend"
	"github.com/oyavri/aldim_verdim/internal/shared/kafka"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Wallet Frontend",
	})

	cfg, err := frontend.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config")
	}

	ctx := context.Background()
	dbPool, err := frontend.NewPostgresPool(ctx, cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Error creating database pool")
	}

	kafkaProducer := kafka.NewProducer()

	repository := frontend.NewWalletRepository(dbPool)
	service := frontend.NewWalletService(repository, kafkaProducer)
	handler := frontend.NewWalletHandler(service)

	app.Get("/", handler.GetWallets)
	app.Post("/", handler.PostEvents)

	app.Listen(fmt.Sprintf("%s:%s", cfg.Hostname, cfg.Hostname))
}
