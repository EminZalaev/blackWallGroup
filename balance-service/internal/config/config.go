package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	// ServerAddr порт на котором будет сервис
	ApiBindAddr        = "BALANCESERVICE_API_BIND_ADDR"
	ApiBindAddrDefault = ":8080"

	// KafkaBrokers kafka brokers
	KafkaBrokers = "BALANCESERVICE_KAFKA_BROKERS"

	KafkaTransferConsumerGroup        = "BALANCESERVICE_KAFKA_TRANSFER_CONSUMER_GROUP"
	KafkaTransferConsumerGroupDefault = "transfer_group"

	KafkaTransferTopic        = "BALANCESERVICE_KAFKA_TOPIC"
	KafkaTransferTopicDefault = "transfer_topic"

	PostgresUser = "BALANCESERVICE_POSTGRES_USER"

	PostgresPassword = "BALANCESERVICE_POSTGRES_PASSWORD"

	PostgresDB = "BALANCESERVICE_POSTGRES_DB"

	PostgresHost = "BALANCESERVICE_POSTGRES_HOST"

	PostgresPort = "BALANCESERVICE_POSTGRES_PORT"
)

func SetDefaults() {
	// ставим default значения
	viper.SetDefault(ApiBindAddr, ApiBindAddrDefault)
	viper.SetDefault(KafkaTransferConsumerGroup, KafkaTransferConsumerGroupDefault)
	viper.SetDefault(KafkaTransferTopic, KafkaTransferTopicDefault)
}

func Configure() error {
	SetDefaults()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error read config env file: %w", err)
	}
	return nil
}
