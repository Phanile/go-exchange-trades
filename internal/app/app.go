package app

import (
	"github.com/Phanile/go-exchange-trades/internal/app/grpc"
	"github.com/Phanile/go-exchange-trades/internal/app/kafka"
	"github.com/Phanile/go-exchange-trades/internal/config"
	"github.com/Phanile/go-exchange-trades/internal/services/trades"
	"log/slog"
)

type App struct {
	GRPCApp  *grpc.App
	KafkaApp *kafka.App
}

func NewApp(log *slog.Logger, kafkaConfig *config.KafkaConfig, gRPCConfig *config.GRPCConfig) *App {
	tradesService := trades.NewTradesService(log)

	gRPCApp := grpc.NewGRPCApp(log, gRPCConfig, tradesService)
	kafkaApp, errKafka := kafka.NewKafkaApp(log, kafkaConfig)

	if errKafka != nil {
		panic(errKafka)
	}

	return &App{
		GRPCApp:  gRPCApp,
		KafkaApp: kafkaApp,
	}
}
