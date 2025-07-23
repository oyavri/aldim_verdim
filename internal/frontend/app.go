package frontend

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	kafkaProducer := NewKafkaProducer(cfg.Broker, cfg.BrokerTopic)

	repository := NewWalletRepository(dbPool)
	service := NewWalletService(repository, kafkaProducer)
	handler := NewWalletHandler(service)

	app := fiber.New(fiber.Config{
		AppName: "Wallet Frontend",
	})

	app.Use(logger.New())
	app.Get("/healthz", handler.HealthCheck)

	app.Get("/", handler.GetWallets)
	app.Post("/", handler.PostEvents)

	go func() {
		log.Printf("Server is starting to listen requests on %s:%s\n", cfg.Hostname, cfg.Port)
		app.Listen(fmt.Sprintf("%s:%s", cfg.Hostname, cfg.Port))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	signal.Stop(signalChan)

	log.Println("App is shutting down")
	err = app.Shutdown()
	if err != nil {
		log.Printf("An error occurred when app is shutting down: %w", err)
	}

	log.Println("Producer is closing")
	err = kafkaProducer.Close()
	if err != nil {
		log.Printf("An error occurred when producer is being closed: %w", err)
	}

	log.Println("Database connection is closing")
	dbPool.Close()
}
