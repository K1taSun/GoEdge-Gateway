package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/k1tasun/GoEdge-Gateway/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting MQTT Ingestor connected to %s", cfg.MQTTBroker)

	// TODO: Initialize MQTT client and subscribe to topics
	// mqttClient := mqtt.NewClient(...)
	// mqttClient.Subscribe(cfg.MQTTTopic, handler)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down MQTT Ingestor...")
}
