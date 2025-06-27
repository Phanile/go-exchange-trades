package trades

import (
	"context"
	tradev1 "github.com/Phanile/go-exchange-protos/generated/go/trades"
	"github.com/Phanile/go-exchange-trades/internal/core"
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
	"log/slog"
)

type Trades struct {
	log               *slog.Logger
	orderProvider     OrderProvider
	tradeProvider     TradeProvider
	orderBookProvider OrderBookProvider
}

func NewTradesService(log *slog.Logger, orderProvider OrderProvider, tradeProvider TradeProvider, orderBookProvider OrderBookProvider) *Trades {
	return &Trades{
		log:               log,
		orderProvider:     orderProvider,
		tradeProvider:     tradeProvider,
		orderBookProvider: orderBookProvider,
	}
}

type OrderProvider interface {
	MatchOrder(order *core.OrderItem)
	GetOrdersByPair(firstCoinId, secondCoinId, orderSideId int64) ([]*models.Order, error)
}

type TradeProvider interface {
	SaveTrade(ctx context.Context, buyOrderId, sellOrderId, amount, price, timestamp int64) (int64, error)
}

type OrderBookProvider interface {
	GetOrderBook() ([]*models.OrderBookEntry, []*models.OrderBookEntry, error)
}

func (t *Trades) CreateOrder(ctx context.Context, req *tradev1.CreateOrderRequest) (*tradev1.CreateOrderResponse, error) {
	return nil, nil
}

func (t *Trades) GetOrderBook(ctx context.Context, resp *tradev1.GetOrderBookRequest) (*tradev1.GetOrderBookResponse, error) {
	return nil, nil
}
