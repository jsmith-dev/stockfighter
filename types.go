package stockfighter

import (
	"fmt"
	"time"
)

const OrderDirectionBuy = "buy"
const OrderDirectionSell = "sell"

const OrderTypeLimit = "limit"
const OrderTypeMarket = "market"
const OrderTypeFillOrKill = "fill-or-kill"
const OrderTypeImmediateOrCancel = "immediate-or-cancel"

type StockInfo struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

func (s StockInfo) String() string {
	return fmt.Sprintf("%v (%v)", s.Symbol, s.Name)
}

type StockQuote struct {
	BidPrice      uint64    `json:"bid"`
	BidSize       uint64    `json:"bidSize"`
	BidDepth      uint64    `json:"bidDepth"`
	AskPrice      uint64    `json:"ask"`
	AskSize       uint64    `json:"askSize"`
	AskDepth      uint64    `json:"askDepth"`
	LastPrice     uint64    `json:"last"`
	LastSize      uint64    `json:"lastSize"`
	LastTradeTime time.Time `json:"lastTrade"`
	QuoteTime     time.Time `json:"quoteTime"`
}

type OrderbookEntry struct {
	Price    uint64 `json:"price"`
	Quantity uint64 `json:"qty"`
	IsBuy    bool   `json:"isBuy"`
}

func (oe OrderbookEntry) String() string {
	if oe.IsBuy {
		return fmt.Sprintf("BUY  $%.2f x %v", float64(oe.Price)/100.0, oe.Quantity)
	}

	return fmt.Sprintf("SELL $%.2f x %v", float64(oe.Price)/100.0, oe.Quantity)
}

type Orderbook struct {
	Bids      []OrderbookEntry `json:"bids"`
	Asks      []OrderbookEntry `json:"asks"`
	Timestamp time.Time        `json:"ts"`
}

type OrderFillInfo struct {
	Price     uint64    `json:"price"`
	Quantity  uint64    `json:"qty"`
	Timestamp time.Time `json:"ts"`
}

type OrderStatus struct {
	Direction        string          `json:"direction"`
	OriginalQuantity uint64          `json:"originalQty"`
	Quantity         uint64          `json:"qty"`
	Price            uint64          `json:"price"`
	OrderType        string          `json:"type"`
	OrderID          int64           `json:"id"`
	Account          string          `json:"account"`
	Timestamp        time.Time       `json:"ts"`
	Fills            []OrderFillInfo `json:"fills"`
	TotalFilled      uint64          `json:"totalFilled"`
	Open             bool            `json:"open"`
}
