package kafka

import (
	"github.com/Phanile/go-exchange-trades/internal/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

type App struct {
	log          *slog.Logger
	consumer     *kafka.Consumer
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
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": "earliest",
	}

	consumer, errConsumer := kafka.NewConsumer(&conf)

	if errConsumer != nil {
		return nil, errConsumer
	}

	errSubscribe := consumer.SubscribeTopics(kafkaConfig.Topics, nil)

	if errSubscribe != nil {
		return nil, errSubscribe
	}

	return &App{
		log:          log,
		consumer:     consumer,
		topics:       kafkaConfig.Topics,
		port:         kafkaConfig.Port,
		workersCount: kafkaConfig.WorkersCount,
		stopCh:       make(chan struct{}),
	}, nil
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
}
