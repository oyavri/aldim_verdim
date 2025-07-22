package frontend

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/oyavri/aldim_verdim/pkg/entity"
)

type Service interface {
	GetWallets(context.Context)
	SendEvents(context.Context)
}

type WalletService struct {
	repo     *WalletRepo
	producer *kafka.Producer
	topic    string
}

func NewWalletService(repo *WalletRepo, producer *kafka.Producer, topic string) *WalletService {
	return &WalletService{
		repo:     repo,
		producer: producer,
		topic:    topic,
	}
}

func (s *WalletService) GetWallets(c context.Context) ([]entity.Wallet, error) {
	wallets, err := s.repo.GetWallets(c)

	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *WalletService) SendTransaction(c context.Context, action string, payload []byte) error {
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
		Key:            []byte(action), // what should key be?
		Value:          payload,
	}, nil) // Throw away any report from the delivery

	if err != nil {
		return err
	}

	return nil
}
