package models

type OrderType int
type OrderSide int

const (
	BUY OrderSide = iota
	SELL
)

const (
	OrderTypeMarket OrderType = iota
	OrderTypeLimit
)

type Order struct {
	Id            int64
	UserId        int64
	SendCoinId    int64
	ReceiveCoinId int64
	OrderType     OrderType
	OrderSide     OrderSide
	Amount        int64
	Price         int64
	Timestamp     int64
}

type OrderBookEntry struct {
	Price  int64
	Amount int64
}
