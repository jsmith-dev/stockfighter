package stockfighter

import (
	"fmt"
	"time"
)

// Order directions.
const (
	OrderDirectionBuy  = "buy"
	OrderDirectionSell = "sell"
)

// Order types (https://starfighter.readme.io/docs/place-new-order#order-types).
const (
	OrderTypeLimit             = "limit"
	OrderTypeMarket            = "market"
	OrderTypeFillOrKill        = "fill-or-kill"
	OrderTypeImmediateOrCancel = "immediate-or-cancel"
)

// A StockInfo represents a stock symbol and its name.
type StockInfo struct {
	// Stock symbol
	Symbol string `json:"symbol"`

	// Stock name
	Name string `json:"name"`
}

func (s StockInfo) String() string {
	return fmt.Sprintf("%v (%v)", s.Symbol, s.Name)
}

// A StockQuote represents a stock quote.
type StockQuote struct {
	// Bid best price, size, and depth
	BidPrice uint64 `json:"bid"`
	BidSize  uint64 `json:"bidSize"`
	BidDepth uint64 `json:"bidDepth"`

	// Ask best price, size, and depth
	AskPrice uint64 `json:"ask"`
	AskSize  uint64 `json:"askSize"`
	AskDepth uint64 `json:"askDepth"`

	// Last trade price, size, and timestamp
	LastPrice     uint64    `json:"last"`
	LastSize      uint64    `json:"lastSize"`
	LastTradeTime time.Time `json:"lastTrade"`

	// Quote update time
	QuoteTime time.Time `json:"quoteTime"`
}

// An OrderbookEntry represents an entry in orderbook.
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

// An Orderbook represents an orderbook for a stock.
type Orderbook struct {
	// Bid entries in the orderbook
	Bids []OrderbookEntry `json:"bids"`

	// Ask entries in the orderbook
	Asks []OrderbookEntry `json:"asks"`

	// Timestamp the orderbook was retrieved
	Timestamp time.Time `json:"ts"`
}

// An OrderFillInfo represents an order fill information.
type OrderFillInfo struct {
	Price     uint64    `json:"price"`
	Quantity  uint64    `json:"qty"`
	Timestamp time.Time `json:"ts"`
}

// An OrderStatus represents the status of an open or closed order.
type OrderStatus struct {
	Direction        string          `json:"direction"`
	OriginalQuantity uint64          `json:"originalQty"`
	Quantity         uint64          `json:"qty"`
	Price            uint64          `json:"price"`
	OrderType        string          `json:"orderType"`
	OrderID          int64           `json:"id"`
	Account          string          `json:"account"`
	Timestamp        time.Time       `json:"ts"`
	Fills            []OrderFillInfo `json:"fills"`
	TotalFilled      uint64          `json:"totalFilled"`
	Open             bool            `json:"open"`
}
