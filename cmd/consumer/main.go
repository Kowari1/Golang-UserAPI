package main

import (
	"context"
	"userapi/internal/config"
	"userapi/internal/kafka"
	"userapi/internal/logger"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	broker, err := config.GetKafkaBroker()
	if err != nil {
		logger.Log.Fatal("failed to load broker", zap.Error(err))
	}

	topic, err := config.GetKafkaTopic()
	if err != nil {
		logger.Log.Fatal("failed to load topic", zap.Error(err))
	}

	consumer := kafka.NewConsumer(broker, topic)

	logger.Log.Info("Kafka consumer started...")

	if err := consumer.Start(context.Background()); err != nil {
		logger.Log.Fatal("consumer failed", zap.Error(err))
	}
}
