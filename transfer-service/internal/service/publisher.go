package service

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"

	"transfer-service/internal/models"
)

type TransferPublisher struct {
	publisher *kafka.Publisher
	topic     string
}

func NewTransferPublisher(brokers []string, topic string) (*TransferPublisher, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(true, true),
	)
	if err != nil {
		return nil, err
	}

	return &TransferPublisher{
		publisher: publisher,
		topic:     topic,
	}, nil
}

func (p *TransferPublisher) SendTransfer(transfer *models.KafkaTransfer) error {
	payload, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	return p.PublishTo(p.topic, payload)
}

func (p *TransferPublisher) PublishTo(topic string, payload []byte) error {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	err := p.publisher.Publish(topic, msg)
	return err
}
