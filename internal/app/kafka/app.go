package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Phanile/go-exchange-trades/internal/config"
	"github.com/Phanile/go-exchange-trades/internal/core"
	"github.com/Phanile/go-exchange-trades/internal/storage"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"strconv"
	"sync"
)

type App struct {
	log          *slog.Logger
	consumer     *kafka.Consumer
	producer     *kafka.Producer
	topics       []string
	port         int
	workersCount int
	stopCh       chan struct{}
}

func NewKafkaApp(log *slog.Logger, kafkaConfig *config.KafkaConfig) (*App, error) {
	const op = "kafka.NewKafkaApp"

	log.With(
		slog.String("op", op),
	)

	conf := kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.BootstrapServers,
		"group.id":          kafkaConfig.GroupID,
		"auto.offset.reset": kafkaConfig.AutoOffsetReset,
	}

	consumer, errConsumer := kafka.NewConsumer(&conf)

	if errConsumer != nil {
		return nil, errConsumer
	}

	errSubscribe := consumer.SubscribeTopics(kafkaConfig.Topics, nil)

	if errSubscribe != nil {
		return nil, errSubscribe
	}

	producer, errProducer := kafka.NewProducer(&conf)

	if errProducer != nil {
		return nil, errProducer
	}

	return &App{
		log:          log,
		consumer:     consumer,
		producer:     producer,
		topics:       kafkaConfig.Topics,
		port:         kafkaConfig.Port,
		workersCount: kafkaConfig.WorkersCount,
		stopCh:       make(chan struct{}),
	}, nil
}

func (a *App) GetProducer() *kafka.Producer {
	return a.producer
}

func (a *App) Run(handler func(msg *kafka.Message) error) {
	const op = "kafkaApp.Run"

	a.log.With(
		slog.String("op", op),
	)

	var wg sync.WaitGroup

	for i := 1; i < a.workersCount+1; i++ {
		wg.Add(1)
		a.log.Info("starting worker with id: " + strconv.Itoa(i))

		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-a.stopCh:
					a.log.Info("stopping worker with id: " + strconv.Itoa(id))
					return
				default:
					message, errRead := a.consumer.ReadMessage(-1)

					if errRead != nil {
						a.log.Error("error reading message:", "workerId", slog.Int("id", id), slog.Any("error", errRead))
						return
					}

					errHandle := handler(message)

					if errHandle != nil {
						a.log.Error("error handle message:", "workerId", slog.Int("id", id), slog.Any("error", errHandle))
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
}

func (a *App) Stop() {
	a.log.Info("kafka app is shutting down")
	close(a.stopCh)
	_ = a.consumer.Close()
	a.producer.Close()
}

func NewKafkaHandler(log *slog.Logger, store *storage.Storage) func(msg *kafka.Message) error {
	return func(msg *kafka.Message) error {
		var trade core.TradeMessage

		err := json.Unmarshal(msg.Value, &trade)
		if err != nil {
			return fmt.Errorf("failed to unmarshal trade message: %w", err)
		}

		_, err = store.SaveTrade(
			context.Background(),
			trade.BuyOrderId,
			trade.SellOrderId,
			trade.Amount,
			trade.Price,
			trade.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("failed to save trade: %w", err)
		}

		log.Info("trade saved", "trade", trade)
		return nil
	}
}
