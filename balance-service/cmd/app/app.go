package app

import (
	"balance-service/internal/config"
	"balance-service/internal/service"
	"balance-service/internal/store"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
)

func Run() error {
	// инициализируем конфиг
	err := config.Configure()
	if err != nil {
		return err
	}

	config.LogVars(config.ApiBindAddr,
		config.KafkaBrokers,
		config.KafkaTransferTopic,
		config.KafkaTransferConsumerGroup)

	store, err := store.NewStorage(viper.GetString(config.PostgresUser),
		viper.GetString(config.PostgresPassword),
		viper.GetString(config.PostgresDB),
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresPort))
	if err != nil {
		return fmt.Errorf("error storage: %w", err)
	}

	defer func() {
		if err := store.CloseDBConnection(); err != nil {
			log.Println("error close db connection: ", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	service.InitSubscribers()

	if err = service.KafkaSubscriberRegistry.Init(); err != nil {
		return fmt.Errorf("error kafka sub init: %w", err)
	}

	go func() {
		if err = service.KafkaSubscriberRegistry.Subscribe(ctx, store); err != nil {
			log.Println("error kafka subscribe: ", err)
			return
		}
	}()

	<-ctx.Done()
	stop()

	err = service.KafkaSubscriberRegistry.Close()
	if err != nil {
		return err
	}

	log.Println("Finish subscribing")
	return nil
}
