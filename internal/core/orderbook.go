package core

import (
	"context"
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
)

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

func (ob *OrderBook) GetOrderBook() ([]*models.OrderBookEntry, []*models.OrderBookEntry, error) {
	return nil, nil, nil
}

func (ob *OrderBook) SaveTrade(ctx context.Context, buyOrderId, sellOrderId, amount, price, timestamp int64) (int64, error) {
	return 0, nil
}

func (ob *OrderBook) GetOrdersByPair(firstCoinId, secondCoinId, orderSideId int64) ([]*models.Order, error) {
	return nil, nil
}

func (ob *OrderBook) MatchOrder(order *OrderItem) {
	switch order.OrderSide {
	case models.SELL:
		for ob.Bids.Len() > 0 {

		}
	case models.BUY:
		if ob.Asks.Len() == 0 {
			return
		}

		top := ob.Asks[0]

		for ob.Asks.Len() > 0 {
			if order.Price > top.Price {

			}
			ask := ob.Asks.Pop().(*OrderItem)

			if ask.Amount == order.Amount && ask.Price == ask.Price {

			}
		}
	default:
		return
	}
}
