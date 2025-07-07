package frontend

import (
	"context"

	"github.com/oyavri/aldim_verdim/internal/shared/kafka"
)

type Service interface {
	GetBalance(context.Context)
	SendEvents(context.Context)
}

type WalletService struct {
	repo          *WalletRepository
	kafkaProducer *kafka.Producer
}

func NewWalletService(repo *WalletRepository, producer *kafka.Producer) *WalletService {
	return &WalletService{
		repo:          repo,
		kafkaProducer: producer,
	}
}

func (s *WalletService) GetBalance(c context.Context) (Wallet, error) {
	var walletId string

	wallet, err := s.repo.GetBalance(c, walletId)

	if err != nil {
		return Wallet{}, err
	}

	return wallet, nil
}

func (s *WalletService) SendEvent(c context.Context, action string, payload []byte) error {
	if err := s.kafkaProducer.Publish(c, action, payload); err != nil {

	}

}
