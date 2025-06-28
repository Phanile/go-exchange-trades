package core

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Phanile/go-exchange-trades/internal/config"
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"sync/atomic"
	"time"
)

var globalOrderId atomic.Int64

type OrderBook struct {
	Asks AsksHeap
	Bids BidsHeap

	producer    *kafka.Producer
	kafkaConfig *config.KafkaConfig
}

func NewOrderBook(producer *kafka.Producer, kafkaConfig *config.KafkaConfig) *OrderBook {
	return &OrderBook{
		Asks:        make(AsksHeap, 0),
		Bids:        make(BidsHeap, 0),
		producer:    producer,
		kafkaConfig: kafkaConfig,
	}
}

func (ob *OrderBook) GetOrderBook() ([]*models.OrderBookEntry, []*models.OrderBookEntry) {
	var asks []*models.OrderBookEntry
	var bids []*models.OrderBookEntry

	for _, ask := range ob.Asks {
		asks = append(asks, &models.OrderBookEntry{
			Price:  ask.Price,
			Amount: ask.Amount,
		})
	}

	for _, bid := range ob.Bids {
		bids = append(bids, &models.OrderBookEntry{
			Price:  bid.Price,
			Amount: bid.Amount,
		})
	}

	return asks, bids
}

func (ob *OrderBook) SaveTrade(ctx context.Context, buyOrderId, sellOrderId, amount, price, timestamp int64) (int64, error) {
	const op = "OrderBook.SaveTrade"

	tradeMessage := &TradeMessage{
		BuyOrderId:  buyOrderId,
		SellOrderId: sellOrderId,
		Amount:      amount,
		Price:       price,
		Timestamp:   timestamp,
	}

	data, err := json.Marshal(tradeMessage)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	errProduce := ob.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &ob.kafkaConfig.Topics[1], // "trades"
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, nil)

	if errProduce != nil {
		return 0, fmt.Errorf("%s: %w", op, errProduce)
	}

	return 1, nil
}

func (ob *OrderBook) GetOrderById(order *OrderItem) *models.Order {
	return nil
}

func (ob *OrderBook) GetOrdersByPair(firstCoinId, secondCoinId, orderSideId int64) ([]*models.Order, error) {
	return nil, nil
}

func (ob *OrderBook) CreateOrder(order *OrderItem) int64 {
	id := globalOrderId.Add(1)
	order.OrderId = id

	switch order.OrderSide {
	case models.SELL:
		heap.Push(&ob.Asks, order)
	case models.BUY:
		heap.Push(&ob.Bids, order)
	default:
		return 0
	}

	return id
}

func (ob *OrderBook) MatchOrder(order *OrderItem) {
	switch order.OrderSide {
	case models.SELL:
		for ob.Bids.Len() > 0 {
			top := ob.Bids.Peek()

			if top.Price == order.Price && top.Amount == order.Amount {
				ob.Bids.Pop()
				_, _ = ob.SaveTrade(context.Background(), top.OrderId, order.OrderId, order.Amount, order.Price, time.Now().UnixNano())
				return
			} else {
				break
			}
		}
	case models.BUY:
		for ob.Asks.Len() > 0 {
			top := ob.Asks.Peek()

			if top.Price == order.Price && top.Amount == order.Amount {
				ob.Asks.Pop()
				_, _ = ob.SaveTrade(context.Background(), top.OrderId, order.OrderId, order.Amount, order.Price, time.Now().UnixNano())
				return
			} else {
				break
			}
		}
	default:
		return
	}
}
