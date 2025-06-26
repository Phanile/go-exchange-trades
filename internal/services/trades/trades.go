package trades

import (
	"context"
	tradev1 "github.com/Phanile/go-exchange-protos/generated/go/trades"
	"log/slog"
)

type Trades struct {
	log *slog.Logger
}

func NewTradesService(log *slog.Logger) *Trades {
	return &Trades{
		log: log,
	}
}

func (service *Trades) CreateOrder(ctx context.Context, req *tradev1.CreateOrderRequest) (*tradev1.CreateOrderResponse, error) {
	return nil, nil
}

func (service *Trades) GetOrderBook(ctx context.Context, resp *tradev1.GetOrderBookRequest) (*tradev1.GetOrderBookResponse, error) {
	return nil, nil
}
