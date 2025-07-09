package core

import (
	"container/heap"
	"context"
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
	"sync/atomic"
	"time"
)

var globalOrderId atomic.Int64

type OrderBook struct {
	Asks AsksHeap
	Bids BidsHeap
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks: make(AsksHeap, 0),
		Bids: make(BidsHeap, 0),
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

	ob.MatchOrder(order)

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
