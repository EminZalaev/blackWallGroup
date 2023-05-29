package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"transfer-service/internal/config"
	"transfer-service/internal/models"
)

var (
	logger = watermill.NewStdLogger(
		true, // debug
		true, // trace
	)
	marshaler = kafka.DefaultMarshaler{}

	brokers = []string{"127.0.0.1:9092"}
)

// createPublisher is a helper function that creates a Publisher, in this case - the Kafka Publisher.
func createPublisher() message.Publisher {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: marshaler,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return kafkaPublisher
}

func getTestEvent() *models.KafkaTransfer {
	return &models.KafkaTransfer{
		Id:     uuid.NewString(),
		Amount: 213123.123,
	}
}

// simulateEvents produces events that will be later consumed.
func simulateEvents(publisher message.Publisher) {
	for i := 0; i < 1; i++ {
		event := getTestEvent()

		payload, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}
		err = publisher.Publish(config.KafkaTopicDefault, message.NewMessage(
			watermill.NewUUID(), // internal uuid of the message, useful for debugging
			payload,
		))

		if err != nil {
			panic(err)
		}

		fmt.Println("message sent")
	}
}

func TestPublisher(t *testing.T) {
	pub := createPublisher()
	simulateEvents(pub)
}
