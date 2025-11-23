package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/k1tasun/GoEdge-Gateway/internal/models"
)

type Ingestor struct {
	client mqtt.Client
	topic  string
}

func NewIngestor(brokerURL, clientID, topic string) (*Ingestor, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("Connected to MQTT Broker: %s", brokerURL)
	})
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", token.Error())
	}

	return &Ingestor{
		client: client,
		topic:  topic,
	}, nil
}

func (i *Ingestor) Start(handler func(reading *models.SensorReading)) error {
	token := i.client.Subscribe(i.topic, 1, func(c mqtt.Client, m mqtt.Message) {
		var reading models.SensorReading
		if err := json.Unmarshal(m.Payload(), &reading); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}

		// Ensure timestamp is set if missing
		if reading.RecordedAt.IsZero() {
			reading.RecordedAt = time.Now()
		}

		log.Printf("Received reading from %s: %f %s", reading.DeviceID, reading.Value, reading.Unit)
		handler(&reading)
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", token.Error())
	}

	log.Printf("Subscribed to topic: %s", i.topic)
	return nil
}

func (i *Ingestor) Close() {
	i.client.Disconnect(250)
}

