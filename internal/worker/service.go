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
	repo   *WalletRepo
	locker *KeyedLocker
}

func NewEventService(repo *WalletRepo) *EventService {
	return &EventService{
		repo:   repo,
		locker: NewKeyedLocker(),
	}
}

func (s *EventService) HandleEvent(ctx context.Context, event entity.Event) error {
	// Lock walletId for sequential processing
	s.locker.Lock(event.WalletId)
	defer s.locker.Unlock(event.WalletId)

	amountStr := event.ActionAttributes.Amount
	amountParsed, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	switch event.ActionType {
	case "BALANCE_INCREASE":
		err = s.repo.IncreaseBalance(ctx, event.WalletId, amountParsed, event.ActionAttributes.Currency)
	case "BALANCE_DECREASE":
		err = s.repo.DecreaseBalance(ctx, event.WalletId, amountParsed, event.ActionAttributes.Currency)
	}

	if err != nil {
		return err
	}

	return nil
}
