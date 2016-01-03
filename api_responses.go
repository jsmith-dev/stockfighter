package stockfighter

import "time"

type apiRespHeartbeat struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type apiRespStocks struct {
	OK     bool        `json:"ok"`
	Error  string      `json:"error"`
	Stocks []StockInfo `json:"symbols"`
}

type apiRespStockOrderbook struct {
	OK          bool             `json:"ok"`
	Error       string           `json:"error"`
	VenueSymbol string           `json:"venue"`
	StockSymbol string           `json:"symbol"`
	Bids        []OrderbookEntry `json:"bids"`
	Asks        []OrderbookEntry `json:"asks"`
	Timestamp   time.Time        `json:"ts"`
}

type apiRespNewStockOrder struct {
	OK               bool            `json:"ok"`
	Error            string          `json:"error"`
	VenueSymbol      string          `json:"venue"`
	StockSymbol      string          `json:"symbol"`
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

type apiRespStockQuote struct {
	OK            bool      `json:"ok"`
	Error         string    `json:"error"`
	VenueSymbol   string    `json:"venue"`
	StockSymbol   string    `json:"symbol"`
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

type apiRespStockOrderStatus struct {
	OK               bool            `json:"ok"`
	Error            string          `json:"error"`
	VenueSymbol      string          `json:"venue"`
	StockSymbol      string          `json:"symbol"`
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

type apiRespAllOrdersStatus struct {
	OK          bool    `json:"ok"`
	Error       string  `json:"error"`
	VenueSymbol string  `json:"venue"`
	Orders      []Order `json:"orders"`
}
