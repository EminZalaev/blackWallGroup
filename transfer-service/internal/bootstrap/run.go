package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"transfer-service/internal/config"
	"transfer-service/internal/controller"
	"transfer-service/internal/middleware"
	"transfer-service/internal/router"
	"transfer-service/internal/service"
)

func Run() error {
	// инициализируем конфиг
	err := config.Configure()
	if err != nil {
		return err
	}

	config.LogVars(config.ApiBindAddr,
		config.KafkaBrokers,
		config.KafkaTopic)

	root := gin.New()
	root.Use(gin.Recovery())

	// создаем ограничение по запросам
	rateLimiter := middleware.NewClientRateLimiter()

	transferPublisher, err := service.NewTransferPublisher(
		viper.GetStringSlice(config.KafkaBrokers),
		viper.GetString(config.KafkaTopic))
	if err != nil {
		return err
	}

	transferController := controller.NewTransferController(transferPublisher)

	transferRouter := router.NewTransferRouter(transferController)

	transferRouter.InitApiV1Group(root.Group("api/v1"), middleware.ClientRateLimit(rateLimiter))

	srv := &http.Server{
		Addr:    viper.GetString(config.ApiBindAddr),
		Handler: root,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Println("error listen server: ", err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
	return nil
}
