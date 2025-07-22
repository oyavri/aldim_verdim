package frontend

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewProducer(broker string, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: topic,
		},
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, key string, message []byte) error {
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: message,
	})
}

func (p *KafkaProducer) Close() error {
	return p.Writer.Close()
}
