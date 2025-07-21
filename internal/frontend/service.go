package frontend

import (
	"context"

	"github.com/oyavri/aldim_verdim/internal/shared/entity"
)

type Service interface {
	GetWallets(context.Context)
	SendEvents(context.Context)
}

type WalletService struct {
	repo          *WalletRepo
	kafkaProducer *KafkaProducer
}

func NewWalletService(repo *WalletRepo, producer *KafkaProducer) *WalletService {
	return &WalletService{
		repo:          repo,
		kafkaProducer: producer,
	}
}

func (s *WalletService) GetWallets(c context.Context) ([]entity.Wallet, error) {
	wallets, err := s.repo.GetWallets(c)

	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *WalletService) SendEvent(c context.Context, action string, payload []byte) error {
	if err := s.kafkaProducer.Publish(c, action, payload); err != nil {
		return err
	}

	return nil
}
