package config

import (
	"os"
)

type Config struct {
	MQTTBroker  string
	MQUITTopic  string
	DatabaseURL string
	ServerPort  string
}

func Load() *Config {
	return &Config{
		MQTTBroker:  getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQUITTopic:  getEnv("MQTT_TOPIC", "sensors/#"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/goedge?sslmode=disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
