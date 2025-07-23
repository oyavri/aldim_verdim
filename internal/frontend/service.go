package frontend

import (
	"context"
	"encoding/json"

	"github.com/oyavri/aldim_verdim/pkg/entity"
)

type Service interface {
	GetWallets(context.Context)
	SendEvents(context.Context)
}

type WalletService struct {
	repo     *WalletRepo
	producer *KafkaProducer
}

func NewWalletService(repo *WalletRepo, producer *KafkaProducer) *WalletService {
	return &WalletService{
		repo:     repo,
		producer: producer,
	}
}

func (s *WalletService) GetWallets(c context.Context) ([]entity.Wallet, error) {
	wallets, err := s.repo.GetWallets(c)

	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *WalletService) SendTransaction(c context.Context, event entity.Event) error {
	serializedEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = s.producer.Produce(c, []byte(event.AppId), serializedEvent)
	if err != nil {
		return err
	}

	return nil
}
