package worker

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/oyavri/aldim_verdim/pkg/entity"
)

var canceledContextError = errors.New("context cancelled")

type Service interface {
	HandleEvent(ctx context.Context, event entity.Event) error
}

type EventService struct {
	repo   *WalletRepo
	sem    chan struct{} // for maximum goroutine limit
	locker *KeyedLocker
}

func NewEventService(repo *WalletRepo, maxConcurrentGoroutine int) *EventService {
	return &EventService{
		repo:   repo,
		sem:    make(chan struct{}, maxConcurrentGoroutine),
		locker: NewKeyedLocker(),
	}
}

func (s *EventService) HandleEvent(ctx context.Context, event entity.Event) error {
	// Increase current goroutine count
	select {
	case s.sem <- struct{}{}:
		defer func() { <-s.sem }()
	case <-ctx.Done():
		fmt.Printf("Context cancelled before handling following event: %v", event)
		return canceledContextError
	}

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
