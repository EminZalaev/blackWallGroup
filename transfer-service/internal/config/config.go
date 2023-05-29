package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	// ServerAddr порт на котором будет сервис
	ApiBindAddr        = "TRANSFERSERVICE_API_BIND_ADDR"
	ApiBindAddrDefault = ":8080"

	// KafkaBrokers kafka brokers
	KafkaBrokers = "TRANSFERSERVICE_KAFKA_BROKERS"

	KafkaTopic        = "TRANSFERSERVICE_KAFKA_TOPIC"
	KafkaTopicDefault = "transfer_topic"

	// ClientRateRps Ограничение количества запросов клиента на API
	ClientRateRps                = "TRANSFERSERVICE_API_RPC_RATELIMIT"
	ClientRateRPSDefault float64 = 10

	// ClientRateBurst Распределение допустимых из ограничения запросов клиента на API
	ClientRateBurst        = "TRANSFERSERVICE_API_RPC_BURST"
	ClientRateBurstDefault = 2
)

func SetDefaults() {
	// ставим default значения
	viper.SetDefault(ApiBindAddr, ApiBindAddrDefault)
	viper.SetDefault(KafkaTopic, KafkaTopicDefault)
	viper.SetDefault(ClientRateRps, ClientRateRPSDefault)
	viper.SetDefault(ClientRateBurst, ClientRateBurstDefault)
}

func Configure() error {
	SetDefaults()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error read config env file: %w", err)
	}
	return nil
}
