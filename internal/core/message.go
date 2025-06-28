package core

type TradeMessage struct {
	BuyOrderId  int64 `json:"buy_order_id"`
	SellOrderId int64 `json:"sell_order_id"`
	Amount      int64 `json:"amount"`
	Price       int64 `json:"price"`
	Timestamp   int64 `json:"timestamp"`
}
