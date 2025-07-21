package frontend

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oyavri/aldim_verdim/internal/shared/db"
)

func Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("Error loading config")
	}

	ctx := context.Background()
	dbPool, err := db.NewPostgresPool(ctx, cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Error creating database pool")
	}

	kafkaProducer := NewProducer(cfg.Brokers, cfg.Topic)

	repository := NewWalletRepository(dbPool)
	service := NewWalletService(repository, kafkaProducer)
	handler := NewWalletHandler(service)

	app := fiber.New(fiber.Config{
		AppName: "Wallet Frontend",
	})

	app.Get("/", handler.GetWallets)
	app.Post("/", handler.PostEvents)

	app.Listen(fmt.Sprintf("%s:%s", cfg.Hostname, cfg.Hostname))
}
