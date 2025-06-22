package core

import (
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHeap_Len(t *testing.T) {
	orderbook := &OrderBook{
		Asks: make(AsksHeap, 0),
		Bids: make(BidsHeap, 0),
	}

	order := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        10,
		Price:         4,
		Timestamp:     time.Now().UnixNano(),
	}

	orderSecond := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        100,
		Price:         4,
		Timestamp:     time.Now().UnixNano(),
	}

	NewOrderItem(&orderbook.Bids, order)
	NewOrderItem(&orderbook.Bids, orderSecond)

	assert.Equal(t, 2, orderbook.Bids.Len())
}

func TestHeap_Less(t *testing.T) {
	orderbook := &OrderBook{
		Asks: make(AsksHeap, 0),
		Bids: make(BidsHeap, 0),
	}

	order := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        10,
		Price:         4,
		Timestamp:     time.Now().UnixNano(),
	}

	orderSecond := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        100,
		Price:         5,
		Timestamp:     time.Now().UnixNano(),
	}

	NewOrderItem(&orderbook.Bids, order)
	NewOrderItem(&orderbook.Bids, orderSecond)

	assert.True(t, orderbook.Bids.Less(0, 1))
	assert.False(t, orderbook.Bids.Less(1, 0))
}

func TestHeap_Swap(t *testing.T) {
	orderbook := &OrderBook{
		Asks: make(AsksHeap, 0),
		Bids: make(BidsHeap, 0),
	}

	order := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        10,
		Price:         4,
		Timestamp:     time.Now().UnixNano(),
	}

	orderSecond := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        100,
		Price:         5,
		Timestamp:     time.Now().UnixNano(),
	}

	orderItem := NewOrderItem(&orderbook.Bids, order)
	orderSecondItem := NewOrderItem(&orderbook.Bids, orderSecond)

	orderbook.Bids.Swap(0, 1)
	assert.Equal(t, orderItem, orderbook.Bids[1])
	assert.Equal(t, orderSecondItem, orderbook.Bids[0])
}

func TestHeap_Pop(t *testing.T) {
	orderbook := &OrderBook{
		Asks: make(AsksHeap, 0),
		Bids: make(BidsHeap, 0),
	}

	order := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        10,
		Price:         4,
		Timestamp:     time.Now().UnixNano(),
	}

	orderSecond := &models.Order{
		UserId:        1,
		SendCoinId:    1,
		ReceiveCoinId: 2,
		OrderType:     1,
		OrderSide:     1,
		Amount:        100,
		Price:         5,
		Timestamp:     time.Now().UnixNano(),
	}

	orderItem := NewOrderItem(&orderbook.Bids, order)
	orderSecondItem := NewOrderItem(&orderbook.Bids, orderSecond)

	assert.Equal(t, orderSecondItem, orderbook.Bids.Pop())
	assert.Equal(t, orderItem, orderbook.Bids.Pop())
}
