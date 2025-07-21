package frontend

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
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

	kafkaProducer := NewProducer(cfg.Brokers, cfg.Topic)

	repository := NewWalletRepository(dbPool)
	service := NewWalletService(repository, kafkaProducer)
	handler := NewWalletHandler(service)

	app := fiber.New(fiber.Config{
		AppName: "Wallet Frontend",
	})

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

	app.Shutdown()
	log.Println("App is successfully shut down")

	dbPool.Close()
	log.Println("Database connection is successfully closed")
}
