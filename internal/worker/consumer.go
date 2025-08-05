package worker

import (
	"context"
	"encoding/json"

	"github.com/oyavri/aldim_verdim/pkg/entity"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	r *kafka.Reader
}

func NewKafkaConsumer(groupId string, brokers []string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			GroupID:  groupId,
			Brokers:  brokers,
			Topic:    topic,
			MaxBytes: 10e6,
		}),
	}
}

func (kc *KafkaConsumer) Consume(ctx context.Context) (entity.Event, error) {
	message, err := kc.r.ReadMessage(ctx)
	if err != nil {
		return entity.Event{}, err
	}

	var event entity.Event

	err = json.Unmarshal(message.Value, &event)
	if err != nil {
		return entity.Event{}, err
	}

	return event, nil
}

func (kc *KafkaConsumer) ConsumeBatch(ctx context.Context, batchSize int) ([]entity.Event, error) {
	var messages []entity.Event

	for range batchSize {
		message, err := kc.Consume(ctx)
		if err != nil {
			return messages, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (kc *KafkaConsumer) Close() error {
	return kc.r.Close()
}
