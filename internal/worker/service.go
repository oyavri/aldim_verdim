package worker

import (
	"context"
	"log"
	"sort"
	"strconv"

	"github.com/oyavri/aldim_verdim/pkg/entity"
)

type Service interface {
	HandleEvent(ctx context.Context, event entity.Event) error
}

type EventService struct {
	repo *WalletRepo
}

func NewEventService(repo *WalletRepo) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) HandleEvents(ctx context.Context, events []entity.Event) error {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time < events[j].Time
	})

	for _, event := range events {
		amountStr := event.ActionAttributes.Amount
		amountParsed, err := strconv.ParseFloat(amountStr, 64)

		if err != nil {
			return err
		}

		log.Printf("Handling event: %v", event)

		switch event.ActionType {
		case "BALANCE_INCREASE":
			err = s.repo.IncreaseBalance(ctx, event.WalletId, amountParsed, event.ActionAttributes.Currency)
		case "BALANCE_DECREASE":
			err = s.repo.DecreaseBalance(ctx, event.WalletId, amountParsed, event.ActionAttributes.Currency)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
