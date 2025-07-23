package frontend

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

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

	// Sort the sent events by time to send them in a FIFO queue on Kafka
	sort.Slice(request.Events, func(i, j int) bool {
		return request.Events[i].Time < request.Events[j].Time
	})

	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(request.Events))
	)

	for _, e := range request.Events {
		event := e
		wg.Add(1)

		go func(event entity.Event) {
			defer wg.Done()

			amountStr := event.ActionAttributes.Amount
			_, err := strconv.ParseFloat(amountStr, 64) // Still pass the amount as string to be consumed

			if err != nil {
				errChan <- fmt.Errorf("invalid amount parameter: %w", err)
				return
			}

			if err := h.service.SendTransaction(c.Context(), event); err != nil {
				errChan <- fmt.Errorf("failed to publish to Kafka: %w", err)
				return
			}
		}(event)

		wg.Wait()
		close(errChan)

		if len(errChan) > 0 {
			var allErrors []string
			for err := range errChan {
				allErrors = append(allErrors, err.Error())
			}

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": allErrors})
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

func (h *WalletHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "healthy"})
}
