package stockfighter

import "fmt"

type ErrorAPITimeout struct{}

func (e *ErrorAPITimeout) Error() string {
	return "API time out"
}

type ErrorUnauthorized struct{}

func (e *ErrorUnauthorized) Error() string {
	return "Not authorized"
}

type ErrorVenueNotFound struct {
	VenueSymbol string
}

func (e *ErrorVenueNotFound) Error() string {
	return "Venue not found: " + e.VenueSymbol
}

type ErrorStockNotFound struct {
	VenueSymbol string
	StockSymbol string
}

func (e *ErrorStockNotFound) Error() string {
	return fmt.Sprintf("Stock not found: %v (venue: %v)", e.StockSymbol, e.VenueSymbol)
}
