package frontend

import (
	"log"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/oyavri/aldim_verdim/pkg/dto"
	"github.com/oyavri/aldim_verdim/pkg/entity"
)

type Handler interface {
	GetWallets(c *fiber.Ctx) error
	PostEvents(c *fiber.Ctx) error
	HealthCheck(c *fiber.Ctx) error
}

type WalletHandler struct {
	service *WalletService
}

func NewWalletHandler(service *WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) PostEvents(c *fiber.Ctx) error {
	c.Accepts("application/json")
	c.AcceptsCharsets("utf-8")

	var request dto.EventRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	sort.Slice(request.Events, func(i, j int) bool {
		return request.Events[i].Time.Before(request.Events[j].Time)
	})

	var failedToParse []entity.Event
	var failedToPublish []entity.Event

	for _, event := range request.Events {
		amountStr := event.ActionAttributes.Amount
		_, err := strconv.ParseFloat(amountStr, 64) // Still pass the amount as string to be consumed

		if err != nil {
			log.Printf("failed to parse float: %v", err)
			failedToParse = append(failedToParse, event)
		}

		if err := h.service.SendTransaction(c.Context(), event); err != nil {
			log.Printf("failed to publish to Kafka: %v", err)
			failedToPublish = append(failedToPublish, event)
		}
	}

	if len(failedToParse) > 0 && len(failedToPublish) > 0 {
		return c.Status(fiber.StatusMultiStatus).JSON(
			fiber.Map{
				"error": fiber.Map{
					"failedToParse":   failedToParse,
					"failedToPublish": failedToPublish,
				},
			})
	}

	if len(failedToParse) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": failedToParse})
	}

	if len(failedToPublish) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": failedToPublish})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "all events have successfully been published"})
}

func (h *WalletHandler) GetWallets(c *fiber.Ctx) error {
	wallets, err := h.service.GetWallets(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch wallets"})
	}

	return c.Status(fiber.StatusOK).JSON(
		dto.WalletResponse{
			Wallets: wallets,
		})
}

func (h *WalletHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "healthy"})
}
