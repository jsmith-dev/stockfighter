package stockfighter

import (
	"testing"

	"fmt"

	"os"
	"strings"

	"github.com/stretchr/testify/assert"
)

const (
	testVenue   = "TESTEX"
	testStock   = "FOOBAR"
	testAccount = "EXB123456"

	testPrice    = uint64(5264)
	testQuantity = uint64(4625)

	testVenueNE = "NOEXIST"
	testStockNE = "NOEXIST"
)

var (
	testApiKey   = ""
	testApiKeyNE = "INVALID_API_KEY"
)

func TestPing(t *testing.T) {
	client := NewClient(testApiKey)

	assert.Nil(t, client.Ping())

	assert.Nil(t, client.PingVenue(testVenue))

	// 404: venue not found
	err := client.PingVenue(testVenueNE)
	_, ok := err.(*ErrorVenueNotFound)
	assert.True(t, ok)
}

func TestListStocks(t *testing.T) {
	client := NewClient(testApiKey)

	stocks, err := client.ListStocks(testVenue)
	assert.Nil(t, err)
	var testStockFound = false
	for _, s := range stocks {
		if s.Symbol == testStock {
			testStockFound = true
			break
		}
	}
	assert.True(t, testStockFound)

	// 404: venue not found
	_, err = client.ListStocks(testVenueNE)
	_, ok := err.(*ErrorVenueNotFound)
	assert.True(t, ok)
}

func TestGetOrderbook(t *testing.T) {
	client := NewClient(testApiKey)

	orderbook, err := client.GetOrderbook(testVenue, testStock)
	assert.Nil(t, err)
	assert.NotNil(t, orderbook)
	assert.NotNil(t, orderbook.Bids)
	assert.NotNil(t, orderbook.Asks)

	// 404: venue not found
	_, err = client.GetOrderbook(testVenueNE, testStock)
	_, ok := err.(*ErrorVenueNotFound)
	assert.True(t, ok)

	// 500: "No venue exists with the symbol XXXX"
	_, err = client.GetOrderbook(testVenue, testStockNE)
	assert.NotNil(t, err)
}

func TestGetAllOrders(t *testing.T) {
	client := NewClient(testApiKey)

	orders, err := client.GetAllOrders(testVenue, testAccount)
	assert.Nil(t, err)
	assert.NotNil(t, orders)

	// 404: venue not found
	_, err = client.GetAllOrders(testVenueNE, testStock)
	_, ok := err.(*ErrorVenueNotFound)
	assert.True(t, ok)

	// 401: unauthorized
	clientNE := NewClient(testApiKeyNE)
	_, err = clientNE.GetAllOrders(testVenueNE, testStock)
	_, ok = err.(*ErrorUnauthorized)
	assert.True(t, ok)
}

func TestGetStockOrders(t *testing.T) {
	client := NewClient(testApiKey)

	orders, err := client.GetStockOrders(testVenue, testAccount, testStock)
	assert.Nil(t, err)
	assert.NotNil(t, orders)

	// 404: venue not found
	_, err = client.GetStockOrders(testVenueNE, testAccount, testStock)
	_, ok := err.(*ErrorVenueNotFound)
	assert.True(t, ok)

	// 401: unauthorized
	clientNE := NewClient(testApiKeyNE)
	_, err = clientNE.GetStockOrders(testVenueNE, testAccount, testStock)
	_, ok = err.(*ErrorUnauthorized)
	assert.True(t, ok)
}

func TestGetQuote(t *testing.T) {
	client := NewClient(testApiKey)

	quote, err := client.GetQuote(testVenue, testStock)
	assert.Nil(t, err)
	assert.NotNil(t, quote)

	// 404: venue or stock not found (this API returns 404 for both errors)
	_, err = client.GetQuote(testVenueNE, testStock)
	_, ok := err.(*ErrorStockNotFound)
	assert.True(t, ok)
	_, err = client.GetQuote(testVenue, testStockNE)
	_, ok = err.(*ErrorStockNotFound)
	assert.True(t, ok)
}

func TestOrderStuffs(t *testing.T) {
	client := NewClient(testApiKey)

	// BUY
	buyOrder, err := client.PlaceOrder(testVenue, testStock, testAccount, testPrice, testQuantity, OrderDirectionBuy, OrderTypeLimit)
	assert.Nil(t, err)
	assert.NotNil(t, buyOrder)
	assert.Equal(t, testAccount, buyOrder.Account)
	assert.Equal(t, OrderDirectionBuy, buyOrder.Direction)
	assert.Equal(t, OrderTypeLimit, buyOrder.OrderType)
	assert.NotNil(t, buyOrder.Fills)
	assert.Equal(t, testQuantity, buyOrder.OriginalQuantity)
	assert.Equal(t, testPrice, buyOrder.Price)
	assert.NotZero(t, buyOrder.OrderID)

	// BUY: check status
	buyOrderStatus, err := client.GetOrder(testVenue, testStock, buyOrder.OrderID)
	assert.Nil(t, err)
	assert.NotNil(t, buyOrderStatus)
	assert.Equal(t, testAccount, buyOrderStatus.Account)
	assert.Equal(t, OrderDirectionBuy, buyOrderStatus.Direction)
	assert.Equal(t, OrderTypeLimit, buyOrderStatus.OrderType)
	assert.NotNil(t, buyOrderStatus.Fills)
	assert.Equal(t, testQuantity, buyOrderStatus.OriginalQuantity)
	assert.Equal(t, testPrice, buyOrderStatus.Price)
	assert.Equal(t, buyOrder.OrderID, buyOrderStatus.OrderID)

	// BUY: cancel
	buyOrderCancel, err := client.CancelOrder(testVenue, testStock, buyOrder.OrderID)
	assert.Nil(t, err)
	assert.NotNil(t, buyOrderCancel)
	assert.Equal(t, testAccount, buyOrderCancel.Account)
	assert.Equal(t, OrderDirectionBuy, buyOrderCancel.Direction)
	assert.Equal(t, OrderTypeLimit, buyOrderCancel.OrderType)
	assert.NotNil(t, buyOrderCancel.Fills)
	assert.Equal(t, testQuantity, buyOrderCancel.OriginalQuantity)
	assert.Equal(t, testPrice, buyOrderCancel.Price)
	assert.Equal(t, buyOrder.OrderID, buyOrderCancel.OrderID)
	assert.False(t, buyOrderCancel.Open)

	// SELL
	sellOrder, err := client.PlaceOrder(testVenue, testStock, testAccount, testPrice, testQuantity, OrderDirectionSell, OrderTypeLimit)
	assert.Nil(t, err)
	assert.NotNil(t, sellOrder)
	assert.Equal(t, testAccount, sellOrder.Account)
	assert.Equal(t, OrderDirectionSell, sellOrder.Direction)
	assert.Equal(t, OrderTypeLimit, sellOrder.OrderType)
	assert.NotNil(t, sellOrder.Fills)
	assert.Equal(t, testQuantity, sellOrder.OriginalQuantity)
	assert.Equal(t, testPrice, sellOrder.Price)
	assert.NotZero(t, sellOrder.OrderID)

	// SELL: check status
	sellOrderStatus, err := client.GetOrder(testVenue, testStock, sellOrder.OrderID)
	assert.Nil(t, err)
	assert.NotNil(t, sellOrderStatus)
	assert.Equal(t, testAccount, sellOrderStatus.Account)
	assert.Equal(t, OrderDirectionSell, sellOrderStatus.Direction)
	assert.Equal(t, OrderTypeLimit, sellOrderStatus.OrderType)
	assert.NotNil(t, sellOrderStatus.Fills)
	assert.Equal(t, testQuantity, sellOrderStatus.OriginalQuantity)
	assert.Equal(t, testPrice, sellOrderStatus.Price)
	assert.Equal(t, sellOrder.OrderID, sellOrderStatus.OrderID)

	// SELL: cancel
	sellOrderCancel, err := client.CancelOrder(testVenue, testStock, sellOrder.OrderID)
	assert.Nil(t, err)
	assert.NotNil(t, sellOrderCancel)
	assert.Equal(t, testAccount, sellOrderCancel.Account)
	assert.Equal(t, OrderDirectionSell, sellOrderCancel.Direction)
	assert.Equal(t, OrderTypeLimit, sellOrderCancel.OrderType)
	assert.NotNil(t, sellOrderCancel.Fills)
	assert.Equal(t, testQuantity, sellOrderCancel.OriginalQuantity)
	assert.Equal(t, testPrice, sellOrderCancel.Price)
	assert.Equal(t, sellOrder.OrderID, sellOrderCancel.OrderID)
	assert.False(t, sellOrderCancel.Open)

	// invalid order
	_, err = client.PlaceOrder(testVenue, testStock, testAccount, testPrice, 0, OrderDirectionBuy, OrderTypeLimit)
	assert.NotNil(t, err)
	_, err = client.PlaceOrder(testVenue, testStock, testAccount, testPrice, testPrice, "invaliddirection", OrderTypeLimit)
	assert.NotNil(t, err)
	_, err = client.PlaceOrder(testVenue, testStock, testAccount, testPrice, testPrice, OrderDirectionSell, "invalidtype")
	assert.NotNil(t, err)

	// checking with invalid params
	_, err = client.GetOrder(testVenue, testStock, 0)
	assert.NotNil(t, err)
	_, err = client.GetOrder(testVenueNE, testStock, sellOrder.OrderID)
	assert.NotNil(t, err)
	_, err = client.GetOrder(testVenue, testStockNE, sellOrder.OrderID)
	assert.NotNil(t, err)

	// cancelling with invalid params
	_, err = client.CancelOrder(testVenue, testStock, 0)
	assert.NotNil(t, err)
	_, err = client.CancelOrder(testVenueNE, testStock, sellOrder.OrderID)
	assert.NotNil(t, err)
	_, err = client.CancelOrder(testVenue, testStockNE, sellOrder.OrderID)
	assert.NotNil(t, err)

	// 401: unauthorized
	clientNE := NewClient(testApiKeyNE)
	_, err = clientNE.PlaceOrder(testVenue, testStock, testAccount, testPrice, testPrice, OrderDirectionBuy, OrderTypeLimit)
	_, ok := err.(*ErrorUnauthorized)
	assert.True(t, ok)
	_, err = clientNE.GetOrder(testVenue, testStock, sellOrder.OrderID)
	_, ok = err.(*ErrorUnauthorized)
	assert.True(t, ok)
	_, err = clientNE.CancelOrder(testVenue, testStock, sellOrder.OrderID)
	_, ok = err.(*ErrorUnauthorized)
	assert.True(t, ok)
}

func init() {
	testApiKey = strings.TrimSpace(os.Getenv("API_KEY"))
	if testApiKey == "" {
		panic(fmt.Errorf("API key ($API_KEY) missing"))
	}
}
