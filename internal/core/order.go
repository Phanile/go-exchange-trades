package core

import (
	"container/heap"
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
)

type OrderItem struct {
	OrderId       int64
	Price         int64
	Amount        int64
	SendCoinId    int64
	ReceiveCoinId int64
	OrderType     models.OrderType
	OrderSide     models.OrderSide
	index         int
}

func NewOrderItem(heap heap.Interface, order *models.Order) *OrderItem {
	item := &OrderItem{
		OrderId:       order.Id,
		Price:         order.Price,
		Amount:        order.Amount,
		SendCoinId:    order.SendCoinId,
		ReceiveCoinId: order.ReceiveCoinId,
		OrderType:     order.OrderType,
		OrderSide:     order.OrderSide,
	}

	heap.Push(item)

	return item
}
