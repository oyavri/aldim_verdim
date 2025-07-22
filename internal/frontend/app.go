package frontend

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oyavri/aldim_verdim/pkg/db"
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

	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": cfg.Broker,
		})
	if err != nil {
		log.Fatal("Failed to create new producer")
	}

	repository := NewWalletRepository(dbPool)
	service := NewWalletService(repository, producer, cfg.BrokerTopic)
	handler := NewWalletHandler(service)

	app := fiber.New(fiber.Config{
		AppName: "Wallet Frontend",
	})

	app.Use(logger.New())

	app.Get("/", handler.GetWallets)
	app.Post("/", handler.PostEvents)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "healthy"})
	})

	go func() {
		log.Printf("Server is starting to listen requests on %s:%s\n", cfg.Hostname, cfg.Port)
		app.Listen(fmt.Sprintf("%s:%s", cfg.Hostname, cfg.Port))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	signal.Stop(signalChan)

	app.Shutdown()
	log.Println("App is successfully shut down")

	dbPool.Close()
	log.Println("Database connection is successfully closed")
}
