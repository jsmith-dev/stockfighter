package stockfighter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const apiBaseURL = "https://api.stockfighter.io/ob/api"

// Client represents a client object you can use Stockfighter APIs.
//
// You can create a new Client using NewClient function.
type Client struct {
	apiKey     string
	httpClient http.Client
}

// NewClient creates a new Client using your API key. This never returns nil.
func NewClient(apiKey string) *Client {
	var api = Client{apiKey: apiKey}

	api.httpClient = http.Client{}

	return &api
}

// Ping checks if the API is up.
//
// Ping returns nil if API is running fine. Otherwise it will return an error.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/heartbeat
func (api *Client) Ping() error {
	resp, err := api.httpClient.Get(apiBaseURL + "/heartbeat")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respData apiRespHeartbeat
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return err
	}

	if !respData.OK {
		return errors.New(respData.Error)
	}

	return nil
}

// PingVenue checks if a venue is up.
//
// PingVenue returns nil if the venue is up. Otherwise it will return an error.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/heartbeat
func (api *Client) PingVenue(venue string) error {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	resp, err := api.httpClient.Get(apiBaseURL + "/venues/" + venue + "/heartbeat")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 500: // timeout
		return &ErrorAPITimeout{}
	case 404: // venue not found
		return &ErrorVenueNotFound{VenueSymbol: venue}
	}

	var respData apiRespHeartbeat
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return err
	}

	if !respData.OK {
		return errors.New(respData.Error)
	}

	return nil
}

// ListStocks lists the stocks available for trading on a venue.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks
func (api *Client) ListStocks(venue string) ([]StockInfo, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	resp, err := api.httpClient.Get(apiBaseURL + "/venues/" + venue + "/stocks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case 404: // venue not found
		return nil, &ErrorVenueNotFound{VenueSymbol: venue}
	}

	var respData apiRespStocks
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return respData.Stocks, nil
}

// GetOrderbook returns the orderbook for a particular stock.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock
func (api *Client) GetOrderbook(venue, stock string) (*Orderbook, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	resp, err := api.httpClient.Get(apiBaseURL + "/venues/" + venue + "/stocks/" + stock)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	var respData apiRespStockOrderbook
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return &Orderbook{
		Bids:      respData.Bids,
		Asks:      respData.Asks,
		Timestamp: respData.Timestamp,
	}, nil
}

// PlaceOrder places an order for a stock.
//
// Stockfighter API:
//     POST https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/orders
func (api *Client) PlaceOrder(venue, stock, account string, price, quantity uint64, direction, orderType string) (*OrderStatus, error) {
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
	req, err := http.NewRequest("POST", apiBaseURL+"/venues/"+venue+"/stocks/"+stock+"/orders", reqBody)
	if err != nil {
		return nil, err
	}
	req.Header = map[string][]string{
		"X-Starfighter-Authorization": {api.apiKey},
		"Content-Type":                {"application/json"},
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	var respData apiRespNewStockOrder
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return &OrderStatus{
		Direction:        respData.Direction,
		OriginalQuantity: respData.OriginalQuantity,
		Quantity:         respData.Quantity,
		Price:            respData.Price,
		OrderType:        respData.OrderType,
		OrderID:          respData.OrderID,
		Account:          respData.Account,
		Timestamp:        respData.Timestamp,
		Fills:            respData.Fills,
		TotalFilled:      respData.TotalFilled,
		Open:             respData.Open,
	}, nil
}

// GetQuote returns a quick look at the most recent trade information for a stock.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/quote
func (api *Client) GetQuote(venue, stock string) (*StockQuote, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	resp, err := api.httpClient.Get(apiBaseURL + "/venues/" + venue + "/stocks/" + stock + "/quote")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	var respData apiRespStockQuote
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return &StockQuote{
		BidPrice:      respData.BidPrice,
		BidSize:       respData.BidSize,
		BidDepth:      respData.BidDepth,
		AskPrice:      respData.AskPrice,
		AskSize:       respData.AskSize,
		AskDepth:      respData.AskDepth,
		LastPrice:     respData.LastPrice,
		LastSize:      respData.LastSize,
		LastTradeTime: respData.LastTradeTime,
		QuoteTime:     respData.QuoteTime,
	}, nil
}

// GetOrder returns a status of an existing order.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/orders/:id
func (api *Client) GetOrder(venue, stock string, orderID int64) (*OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	req, err := http.NewRequest("GET", apiBaseURL+"/venues/"+venue+"/stocks/"+stock+"/orders/"+strconv.FormatInt(orderID, 10), nil)
	req.Header = map[string][]string{
		"X-Starfighter-Authorization": {api.apiKey},
	}
	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	}

	var respData apiRespStockOrderStatus
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return &OrderStatus{
		Direction:        respData.Direction,
		OriginalQuantity: respData.OriginalQuantity,
		Quantity:         respData.Quantity,
		Price:            respData.Price,
		OrderType:        respData.OrderType,
		OrderID:          respData.OrderID,
		Account:          respData.Account,
		Timestamp:        respData.Timestamp,
		Fills:            respData.Fills,
		TotalFilled:      respData.TotalFilled,
		Open:             respData.Open,
	}, nil
}

// CancelOrder cancels an order.
//
// Stockfighter API:
//     DELETE https://api.stockfighter.io/ob/api/venues/:venue/stocks/:stock/orders/:order
func (api *Client) CancelOrder(venue, stock string, orderID int64) (*OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	stock = strings.TrimSpace(stock)
	if stock == "" {
		panic(fmt.Errorf("Invalid stock symbol: %v", stock))
	}

	req, err := http.NewRequest("DELETE", apiBaseURL+"/venues/"+venue+"/stocks/"+stock+"/orders/"+strconv.FormatInt(orderID, 10), nil)
	req.Header = map[string][]string{
		"X-Starfighter-Authorization": {api.apiKey},
	}
	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	case 404: // stock not found
		return nil, &ErrorStockNotFound{VenueSymbol: venue, StockSymbol: stock}
	}

	var respData apiRespStockOrderStatus
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return &OrderStatus{
		Direction:        respData.Direction,
		OriginalQuantity: respData.OriginalQuantity,
		Quantity:         respData.Quantity,
		Price:            respData.Price,
		OrderType:        respData.OrderType,
		OrderID:          respData.OrderID,
		Account:          respData.Account,
		Timestamp:        respData.Timestamp,
		Fills:            respData.Fills,
		TotalFilled:      respData.TotalFilled,
		Open:             respData.Open,
	}, nil
}

// GetAllOrders returns status of all stock orders in the venue.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/accounts/:account/orders
func (api *Client) GetAllOrders(venue, account string) ([]OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	account = strings.TrimSpace(account)
	if account == "" {
		panic(fmt.Errorf("Invalid account name: %v", account))
	}

	req, err := http.NewRequest("GET", apiBaseURL+"/venues/"+venue+"/accounts/"+account+"/orders", nil)
	req.Header = map[string][]string{
		"X-Starfighter-Authorization": {api.apiKey},
	}
	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	}

	var respData apiRespAllOrdersStatus
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return respData.Orders, nil
}

// GetStockOrders returns status of all orders for a particular stock in the venue.
//
// Stockfighter API:
//     GET https://api.stockfighter.io/ob/api/venues/:venue/accounts/:account/stocks/:stock/orders
func (api *Client) GetStockOrders(venue, account, stock string) ([]OrderStatus, error) {
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

	req, err := http.NewRequest("GET", apiBaseURL+"/venues/"+venue+"/accounts/"+account+"/stocks/"+stock+"/orders", nil)
	req.Header = map[string][]string{
		"X-Starfighter-Authorization": {api.apiKey},
	}
	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401: // unauthorized
		return nil, &ErrorUnauthorized{}
	}

	var respData apiRespAllOrdersStatus
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, err
	}

	if !respData.OK {
		return nil, errors.New(respData.Error)
	}

	return respData.Orders, nil
}
