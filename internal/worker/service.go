package worker

import (
	"context"
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

func (s *EventService) HandleEvent(ctx context.Context, event entity.Event) error {
	amountStr := event.ActionAttributes.Amount
	amountParsed, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	switch event.ActionType {
	case "BALANCE_INCREASE":
		s.repo.IncreaseBalance(ctx, event.Meta.UserId, event.WalletId, amountParsed, event.ActionAttributes.Currency)
	case "BALANCE_DECREASE":
		s.repo.DecreaseBalance(ctx, event.Meta.UserId, event.WalletId, amountParsed, event.ActionAttributes.Currency)
	}

	return nil
}
