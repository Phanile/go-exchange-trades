package trades

import (
	"context"
	tradev1 "github.com/Phanile/go-exchange-protos/generated/go/trades"
)

type Trades interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrderBook(context.Context, *GetOrderBookRequest) (*GetOrderBookResponse, error)
}

type ServerAPI struct {
	tradev1.UnimplementedTradeServer
	trades Trades
}

func Register() {
	tradev1.RegisterTradeServer()
}
