package stockfighter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Client represents a client object you can use Stockfighter APIs.
//
// You can create a new Client using NewClient function.
type Client struct {
	apiKey     string
	apiBaseURL string
	httpClient http.Client
}

// NewClient creates a new Client using your API key. This never returns nil.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		apiBaseURL: "https://api.stockfighter.io/ob/api",
		httpClient: http.Client{},
	}
}

func (client *Client) getAPIJson(method, apiPath string, reqBody io.Reader, respBody interface{}) (int, error) {
	req, err := http.NewRequest(strings.ToUpper(method), client.apiBaseURL+apiPath, reqBody)
	if err != nil {
		return 0, err
	}
	req.Header = map[string][]string{
		"X-Starfighter-Authorization": {client.apiKey},
		"Content-Type":                {"application/json"},
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return resp.StatusCode, decoder.Decode(respBody)
}

// Ping checks if the API is up.
//
// Ping returns nil if API is running fine. Otherwise it will return an error.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/heartbeat
func (client *Client) Ping() error {
	var resp apiRespHeartbeat
	_, err := client.getAPIJson("GET", "/heartbeat", nil, &resp)
	if err != nil {
		return err
	}

	if !resp.OK {
		return errors.New(resp.Error)
	}

	return nil
}

// PingVenue checks if a venue is up.
//
// PingVenue returns nil if the venue is up. Otherwise it will return an error.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/heartbeat
func (client *Client) PingVenue(venue string) error {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	var resp apiRespHeartbeat
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/heartbeat", nil, &resp)
	switch {
	case err != nil:
		return err
	case status == 500: // timeout
		return &ErrorAPITimeout{}
	case status == 404: // venue not found
		return &ErrorVenueNotFound{VenueSymbol: venue}
	}

	if !resp.OK {
		return errors.New(resp.Error)
	}

	return nil
}

// ListStocks lists the stocks available for trading on a venue.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks
func (client *Client) ListStocks(venue string) ([]StockInfo, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	var resp apiRespStocks
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/stocks", nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case status == 404: // venue not found
		return nil, &ErrorVenueNotFound{VenueSymbol: venue}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return resp.Stocks, nil
}

// GetOrderbook returns the orderbook for a particular stock.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock
func (client *Client) GetOrderbook(venue, stock string) (*Orderbook, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	var resp apiRespStockOrderbook
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/stocks/"+stock, nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case status == 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return &Orderbook{
		Bids:      resp.Bids,
		Asks:      resp.Asks,
		Timestamp: resp.Timestamp,
	}, nil
}

// PlaceOrder places an order for a stock.
//
// Stockfighter API:
//     POST https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/orders
func (client *Client) PlaceOrder(venue, stock, account string, price, quantity uint64, direction, orderType string) (*OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	account = strings.TrimSpace(account)
	if account == "" {
		panic(fmt.Errorf("Invalid account name: %v", account))
	}

	reqBody := strings.NewReader(fmt.Sprintf(`{
			"account": "%s",
			"venue": "%s",
			"stock": "%s",
			"price": %d,
			"qty": %d,
			"direction": "%s",
			"orderType": "%s"
		}`, account, venue, stock, price, quantity, direction, orderType))

	var resp apiRespNewStockOrder
	status, err := client.getAPIJson("POST", "/venues/"+venue+"/stocks/"+stock+"/orders", reqBody, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case status == 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return &OrderStatus{
		Direction:        resp.Direction,
		OriginalQuantity: resp.OriginalQuantity,
		Quantity:         resp.Quantity,
		Price:            resp.Price,
		OrderType:        resp.OrderType,
		OrderID:          resp.OrderID,
		Account:          resp.Account,
		Timestamp:        resp.Timestamp,
		Fills:            resp.Fills,
		TotalFilled:      resp.TotalFilled,
		Open:             resp.Open,
	}, nil
}

// GetQuote returns a quick look at the most recent trade information for a stock.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/quote
func (client *Client) GetQuote(venue, stock string) (*StockQuote, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	var resp apiRespStockQuote
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/stocks/"+stock+"/quote", nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case status == 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return &StockQuote{
		BidPrice:      resp.BidPrice,
		BidSize:       resp.BidSize,
		BidDepth:      resp.BidDepth,
		AskPrice:      resp.AskPrice,
		AskSize:       resp.AskSize,
		AskDepth:      resp.AskDepth,
		LastPrice:     resp.LastPrice,
		LastSize:      resp.LastSize,
		LastTradeTime: resp.LastTradeTime,
		QuoteTime:     resp.QuoteTime,
	}, nil
}

// GetOrder returns a status of an existing order.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/orders/:id
func (client *Client) GetOrder(venue, stock string, orderID int64) (*OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	var resp apiRespStockOrderStatus
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/stocks/"+stock+"/orders/"+strconv.FormatInt(orderID, 10), nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return &OrderStatus{
		Direction:        resp.Direction,
		OriginalQuantity: resp.OriginalQuantity,
		Quantity:         resp.Quantity,
		Price:            resp.Price,
		OrderType:        resp.OrderType,
		OrderID:          resp.OrderID,
		Account:          resp.Account,
		Timestamp:        resp.Timestamp,
		Fills:            resp.Fills,
		TotalFilled:      resp.TotalFilled,
		Open:             resp.Open,
	}, nil
}

// CancelOrder cancels an order.
//
// Stockfighter API:
//     DELETE https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/orders/:order
func (client *Client) CancelOrder(venue, stock string, orderID int64) (*OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	var resp apiRespStockOrderStatus
	status, err := client.getAPIJson("DELETE", "/venues/"+venue+"/stocks/"+stock+"/orders/"+strconv.FormatInt(orderID, 10), nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case status == 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return &OrderStatus{
		Direction:        resp.Direction,
		OriginalQuantity: resp.OriginalQuantity,
		Quantity:         resp.Quantity,
		Price:            resp.Price,
		OrderType:        resp.OrderType,
		OrderID:          resp.OrderID,
		Account:          resp.Account,
		Timestamp:        resp.Timestamp,
		Fills:            resp.Fills,
		TotalFilled:      resp.TotalFilled,
		Open:             resp.Open,
	}, nil
}

// GetAllOrders returns status of all stock orders in the venue.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/accounts/:account/orders
func (client *Client) GetAllOrders(venue, account string) ([]OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	account = strings.TrimSpace(account)
	if account == "" {
		panic(fmt.Errorf("Invalid account name: %v", account))
	}

	var resp apiRespAllOrdersStatus
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/accounts/"+account+"/orders", nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return resp.Orders, nil
}

// GetStockOrders returns status of all orders for a particular stock in the venue.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/accounts/:account/stocks/:stock/orders
func (client *Client) GetStockOrders(venue, account, stock string) ([]OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	account = strings.TrimSpace(account)
	if account == "" {
		panic(fmt.Errorf("Invalid account name: %v", account))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	var resp apiRespAllOrdersStatus
	status, err := client.getAPIJson("GET", "/venues/"+venue+"/accounts/"+account+"/stocks/"+stock+"/orders", nil, &resp)
	switch {
	case err != nil:
		return nil, err
	case status == 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	}

	if !resp.OK {
		return nil, errors.New(resp.Error)
	}

	return resp.Orders, nil
}
