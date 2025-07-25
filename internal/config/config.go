package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	return godotenv.Load()
}

func GetJwtKey() ([]byte, error) {
	key := os.Getenv("JWT_KEY")

	if key == "" {
		return nil, fmt.Errorf("JWT_KEY must be set in .env or environment")
	}

	return []byte(key), nil
}

func GetDBDsn() (string, error) {
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		return "", fmt.Errorf("DB_DSN must be set")
	}

	return dsn, nil
}

func GetPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return ":" + port
}

func GetJwtExpiration() time.Duration {
	s := os.Getenv("JWT_EXP_MINUTES")

	if s == "" {
		return 24 * time.Hour
	}

	mins, err := strconv.Atoi(s)

	if err != nil {
		return 24 * time.Hour
	}

	return time.Duration(mins) * time.Minute
}

func GetKafkaBroker() (string, error) {
	broker := os.Getenv("KAFKA_BROKER")

	if broker == "" {
		return "", fmt.Errorf("KAFKA_BROKER must be set")
	}

	return broker, nil
}

func GetKafkaTopic() (string, error) {
	topic := os.Getenv("KAFKA_TOPIC")

	if topic == "" {
		return "", fmt.Errorf("KAFKA_TOPIC must be set")
	}

	return topic, nil
}
