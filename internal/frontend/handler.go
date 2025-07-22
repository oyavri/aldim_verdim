package frontend

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/oyavri/aldim_verdim/pkg/dto"
)

type Handler interface {
	GetWallets(*fiber.Ctx)
	PostEvents(*fiber.Ctx)
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

	for _, event := range request.Events {
		payload, err := json.Marshal(event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "serialization failed"})
		}

		if err := h.service.SendTransaction(c.Context(), event.ActionType, payload); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to publish to Kafka"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successfully published to Kafka"})
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
