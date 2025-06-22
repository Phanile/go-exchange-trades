package core

import "github.com/Phanile/go-exchange-trades/internal/domain/models"

type OrderBook struct {
	Asks AsksHeap
	Bids BidsHeap
}

func matchOrder(ob *OrderBook, order *OrderItem) {
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
