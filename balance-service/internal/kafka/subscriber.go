package kafka

import (
	"balance-service/internal/config"
	"balance-service/internal/models"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spf13/viper"
	"log"
)

var (
	KafkaSubscriberConfig   *sarama.Config
	KafkaSubscriberRegistry KafkaSubscriber
)

func InitSubscribers() {
	KafkaSubscriberRegistry = KafkaSubscriber{
		Topic:         viper.GetString(config.KafkaTransferTopic),
		ConsumerGroup: viper.GetString(config.KafkaTransferConsumerGroup),
		Callback:      SubscribeCronResume,
	}
}

type KafkaSubscriber struct {
	ConsumerGroup string
	Topic         string
	Callback      func(ctx context.Context, messages <-chan *message.Message)
	messages      <-chan *message.Message
	subscriber    message.Subscriber
}

func (ks *KafkaSubscriber) Init() error {
	if ks.subscriber != nil {
		return nil
	}

	conf := kafka.SubscriberConfig{
		Brokers:       viper.GetStringSlice(config.KafkaBrokers),
		ConsumerGroup: ks.ConsumerGroup,
		Unmarshaler:   kafka.DefaultMarshaler{},
	}

	logger := watermill.NewStdLogger(false, false)

	sub, err := kafka.NewSubscriber(conf, logger)
	if err != nil {
		return err
	}

	ks.subscriber = sub

	return nil
}

func (ks *KafkaSubscriber) Close() error {
	return ks.subscriber.Close()
}

func (ks *KafkaSubscriber) Subscribe(ctx context.Context) error {
	subMessage, err := ks.subscriber.Subscribe(ctx, ks.Topic)
	if err != nil {
		return err
	}

	ks.messages = subMessage
	ks.Callback(ctx, ks.messages)

	return nil
}

func SubscribeCronResume(ctx context.Context, messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		cronResumeMsg := models.TransferMessage{}
		err := json.Unmarshal(msg.Payload, &cronResumeMsg)
		if err != nil {
			log.Printf("error unmarshall msg from kafka: %v", msg.UUID)
			msg.Ack()
			continue
		}

		go handleCronResumeMsg(ctx, cronResumeMsg)

		msg.Ack()
	}
}
