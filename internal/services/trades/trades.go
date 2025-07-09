package trades

import (
	"context"
	"fmt"
	tradev1 "github.com/Phanile/go-exchange-protos/generated/go/trades"
	"github.com/Phanile/go-exchange-trades/internal/core"
	"github.com/Phanile/go-exchange-trades/internal/domain/models"
	"log/slog"
	"strconv"
	"time"
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
	CreateOrder(order *core.OrderItem) int64
	GetOrderById(order *core.OrderItem) *models.Order
	GetOrdersByPair(firstCoinId, secondCoinId, orderSideId int64) ([]*models.Order, error)
}

type TradeProvider interface {
	SaveTrade(ctx context.Context, buyOrderId, sellOrderId, amount, price, timestamp int64) (int64, error)
}

type OrderBookProvider interface {
	GetOrderBook() ([]*models.OrderBookEntry, []*models.OrderBookEntry)
}

func (t *Trades) CreateOrder(ctx context.Context, req *tradev1.CreateOrderRequest) (*tradev1.CreateOrderResponse, error) {
	const op = "trades.CreateOrder"

	t.log.With(
		slog.String("op", op),
	)

	var price int64

	if req.Price != nil {
		priceParsed, errPrice := strconv.ParseInt(*req.Price, 10, 64)

		if errPrice != nil {
			return nil, fmt.Errorf("%s: %w", op, errPrice)
		}

		price = priceParsed
	} else {
		price = 0
	}

	amountParsed, errAmount := strconv.ParseInt(req.Amount, 10, 64)

	if errAmount != nil {
		return nil, fmt.Errorf("%s: %w", op, errAmount)
	}

	orderId := t.orderProvider.CreateOrder(core.NewOrderItem(&models.Order{
		UserId:        req.UserId,
		SendCoinId:    req.FirstCoinId,
		ReceiveCoinId: req.SecondCoinId,
		OrderType:     models.OrderType(req.Type),
		OrderSide:     models.OrderSide(req.Side),
		Amount:        amountParsed,
		Price:         price,
		Timestamp:     time.Now().UnixNano(),
	}))

	return &tradev1.CreateOrderResponse{
		OrderId: orderId,
	}, nil
}

func (t *Trades) GetOrderBook(ctx context.Context, resp *tradev1.GetOrderBookRequest) (*tradev1.GetOrderBookResponse, error) {
	const op = "trades.GetOrderBook"

	t.log.With(
		slog.String("op", op),
	)

	bids, asks := t.orderBookProvider.GetOrderBook()

	return &tradev1.GetOrderBookResponse{
		Bids: mapOrderBookEntries(bids),
		Asks: mapOrderBookEntries(asks),
	}, nil
}

func mapOrderBookEntries(entries []*models.OrderBookEntry) []*tradev1.OrderBookEntry {
	result := make([]*tradev1.OrderBookEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &tradev1.OrderBookEntry{
			Price:  strconv.FormatInt(entry.Price, 10),
			Amount: strconv.FormatInt(entry.Amount, 10),
		})
	}
	return result
}
