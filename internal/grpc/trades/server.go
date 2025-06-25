package trades

import (
	"context"
	tradev1 "github.com/Phanile/go-exchange-protos/generated/go/trades"
	"google.golang.org/grpc"
)

type Trades interface {
	CreateOrder(ctx context.Context, req *tradev1.CreateOrderRequest) (*tradev1.CreateOrderResponse, error)
	GetOrderBook(ctx context.Context, resp *tradev1.GetOrderBookRequest) (*tradev1.GetOrderBookResponse, error)
}

type ServerAPI struct {
	tradev1.UnimplementedTradeServer
	trades Trades
}

func Register(s *grpc.Server, trades Trades) {
	tradev1.RegisterTradeServer(s, &ServerAPI{
		trades: trades,
	})
}

func (s *ServerAPI) CreateOrder(ctx context.Context, req *tradev1.CreateOrderRequest) (*tradev1.CreateOrderResponse, error) {
	return &tradev1.CreateOrderResponse{
		OrderId: 1,
	}, nil
}

func (s *ServerAPI) GetOrderBook(ctx context.Context, resp *tradev1.GetOrderBookRequest) (*tradev1.GetOrderBookResponse, error) {
	return &tradev1.GetOrderBookResponse{
		Asks: make([]*tradev1.OrderBookEntry, 0),
		Bids: make([]*tradev1.OrderBookEntry, 0),
	}, nil
}
