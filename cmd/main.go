package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/prasannakumar414/somq/config"
	mongodata "github.com/prasannakumar414/somq/data/mongo"
	somqhttp "github.com/prasannakumar414/somq/http"
	"github.com/prasannakumar414/somq/http/handlers"
	"github.com/prasannakumar414/somq/message-queue/kafka"
	"github.com/prasannakumar414/somq/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// ── Config ────────────────────────────────────────────────────────────────
	configPath := getEnv("CONFIG_PATH", "config/config.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatal("Failed to load config", zap.String("path", configPath), zap.Error(err))
	}
	logger.Info("Config loaded", zap.String("path", configPath))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// ── MongoDB ──────────────────────────────────────────────────────────────
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	// ── Kafka ────────────────────────────────────────────────────────────────
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true

	kafkaProducerClient, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, saramaConfig)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Strings("brokers", cfg.Kafka.Brokers), zap.Error(err))
	}
	defer func() {
		if err := kafkaProducerClient.Close(); err != nil {
			logger.Error("Failed to close Kafka producer", zap.Error(err))
		}
	}()
	logger.Info("Connected to Kafka", zap.Strings("brokers", cfg.Kafka.Brokers))

	// ── Wire dependencies ────────────────────────────────────────────────────
	scheduleMessageDAO := mongodata.NewScheduleMessageDAO(mongoClient, logger)
	messageProducer := kafka.NewMessageProducer(logger, kafkaProducerClient)
	schedulerService := services.NewSchedulerService(scheduleMessageDAO, logger, messageProducer)
	messageService := services.NewMessageService(logger, scheduleMessageDAO, messageProducer)

	// ── Scheduler ─────────────────────────────────────────────────────────────
	go schedulerService.Run(ctx)
	logger.Info("Scheduler started")

	messageHandler := handlers.NewMessageHandler(logger, messageService)

	// ── HTTP server ───────────────────────────────────────────────────────────
	server := somqhttp.NewServer(cfg.Server.Port, *logger, messageHandler)
	go func() {
		if err := server.Serve(); err != nil {
			logger.Error("HTTP server stopped", zap.Error(err))
			cancel()
		}
	}()
	logger.Info("HTTP server started", zap.Int("port", cfg.Server.Port))

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	<-ctx.Done()
	logger.Info("Shutting down...")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
