package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/oyavri/aldim_verdim/internal/frontend"
	"github.com/oyavri/aldim_verdim/internal/shared/kafka"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Wallet Frontend",
	})

	kafkaProducer := kafka.NewProducer()

	repository := frontend.NewWalletRepository()
	service := frontend.NewWalletService(repository)
	handler := frontend.NewWalletHandler(kafkaProducer, service)

	app.Get("/", handler.GetWallets)
	app.Post("/", handler.PostEvents)

	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf("%s:%s", hostname, port))
}
