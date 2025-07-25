package kafka

import (
	"context"
	"encoding/json"

	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokerURL, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) SendMessage(ctx context.Context, value interface{}) error {
	bytes, err := json.Marshal(value)

	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: bytes,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
