package app

import (
	"database/sql"
	"github.com/Phanile/go-exchange-trades/internal/app/grpc"
	"github.com/Phanile/go-exchange-trades/internal/app/kafka"
	"github.com/Phanile/go-exchange-trades/internal/config"
	"github.com/Phanile/go-exchange-trades/internal/core"
	"github.com/Phanile/go-exchange-trades/internal/middleware"
	"github.com/Phanile/go-exchange-trades/internal/services/trades"
	"github.com/Phanile/go-exchange-trades/internal/storage"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
)

type App struct {
	GRPCApp  *grpc.App
	KafkaApp *kafka.App
	Storage  *storage.Storage
}

func NewApp(log *slog.Logger, kafkaConfig *config.KafkaConfig, gRPCConfig *config.GRPCConfig) *App {
	postgres, errPostgres := storage.NewPostgresStorage(os.Getenv("PGSQL_CONNECTION_STRING"))

	if errPostgres != nil {
		panic(errPostgres)
	}

	runMigrations(postgres.Connection())

	jwtMiddleware := middleware.NewJWTMiddleware(os.Getenv("JWT_PUBLIC_KEY"))

	kafkaApp, errKafka := kafka.NewKafkaApp(log, kafkaConfig)

	if errKafka != nil {
		panic(errKafka)
	}

	producer := kafkaApp.GetProducer()

	ob := core.NewOrderBook(producer, kafkaConfig)

	tradesService := trades.NewTradesService(log, ob, ob, ob, producer)

	gRPCApp := grpc.NewGRPCApp(log, gRPCConfig, tradesService, jwtMiddleware)

	return &App{
		GRPCApp:  gRPCApp,
		KafkaApp: kafkaApp,
		Storage:  postgres,
	}
}

func runMigrations(db *sql.DB) {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
