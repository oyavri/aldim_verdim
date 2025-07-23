package frontend

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	w *kafka.Writer
}

func NewKafkaProducer(broker string, topic string) *KafkaProducer {
	return &KafkaProducer{
		w: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (kp *KafkaProducer) Produce(ctx context.Context, key []byte, value []byte) error {
	message := &kafka.Message{
		Key:   key,
		Value: value,
	}

	err := kp.w.WriteMessages(ctx, *message)

	if err != nil {
		return err
	}

	return nil
}

func (kp *KafkaProducer) Close() {
	kp.w.Close()
}
