package main

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

//TODO: Каждый GRPC запрос попадает в Kafka который кладет данные в Postgre
//TODO: Создать Kafka app

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "trades_service",
	}

	producer, err := kafka.NewProducer(config)

	if err != nil {
		panic(err)
	}
}
