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

type APIClient struct {
	apiKey     string
	httpClient http.Client
}

func NewAPIClient(apiKey string) *APIClient {
	var api = APIClient{apiKey: apiKey}

	api.httpClient = http.Client{}

	return &api
}

func (api *APIClient) Ping() error {
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

func (api *APIClient) PingVenue(venue string) error {
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

func (api *APIClient) ListStocks(venue string) ([]StockInfo, error) {
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

func (api *APIClient) GetStockOrderbook(venue, stock string) (*Orderbook, error) {
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

func (api *APIClient) PlaceStockOrder(venue, stock, account string, price, quantity uint64, direction, orderType string) (*OrderStatus, error) {
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

func (api *APIClient) GetStockQuote(venue, stock string) (*StockQuote, error) {
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

func (api *APIClient) GetStockOrderStatus(venue, stock string, orderID int64) (*OrderStatus, error) {
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

func (api *APIClient) CancelStockOrder(venue, stock string, orderID int64) (*OrderStatus, error) {
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

func (api *APIClient) GetAllStockAllOrdersStatus(venue, account string) ([]OrderStatus, error) {
	venue = strings.TrimSpace(venue)
	if venue == "" {
		panic(fmt.Errorf("Invalid venue symbol: %v", venue))
	}

	account = strings.TrimSpace(account)
	if account == "" {
		panic(fmt.Errorf("Invalid account name: %v", account))
	}

	resp, err := api.httpClient.Get(apiBaseURL + "/venues/" + venue + "/accounts/" + account + "/orders")
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

func (api *APIClient) GetStockAllOrdersStatus(venue, account, stock string) ([]OrderStatus, error) {
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

	resp, err := api.httpClient.Get(apiBaseURL + "/venues/" + venue + "/accounts/" + account + "/stocks/" + stock + "/orders")
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
