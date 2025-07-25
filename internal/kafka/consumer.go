package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewConsumer(brocker, topic string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brocker},
		Topic:   topic,
		GroupID: "user-consumer-group",
	})

	return &KafkaConsumer{reader: r}
}

func (c KafkaConsumer) Start(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("🎯 Received message: %s\n", string(m.Value))
	}
}
