package app

import (
	"github.com/Phanile/go-exchange-trades/internal/app/grpc"
	"github.com/Phanile/go-exchange-trades/internal/app/kafka"
	"github.com/Phanile/go-exchange-trades/internal/config"
	"github.com/Phanile/go-exchange-trades/internal/middleware"
	"github.com/Phanile/go-exchange-trades/internal/services/trades"
	"log/slog"
	"os"
)

type App struct {
	GRPCApp  *grpc.App
	KafkaApp *kafka.App
}

func NewApp(log *slog.Logger, kafkaConfig *config.KafkaConfig, gRPCConfig *config.GRPCConfig) *App {
	tradesService := trades.NewTradesService(log)
	jwtMiddleware := middleware.NewJWTMiddleware(os.Getenv("JWT_PUBLIC_KEY"))

	gRPCApp := grpc.NewGRPCApp(log, gRPCConfig, tradesService, jwtMiddleware)
	kafkaApp, errKafka := kafka.NewKafkaApp(log, kafkaConfig)

	if errKafka != nil {
		panic(errKafka)
	}

	return &App{
		GRPCApp:  gRPCApp,
		KafkaApp: kafkaApp,
	}
}
