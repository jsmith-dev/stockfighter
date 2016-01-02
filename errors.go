package stockfighter

import "fmt"

// API timeout error.
type ErrorAPITimeout struct{}

func (e *ErrorAPITimeout) Error() string {
	return "API time out"
}

// Unauthorized error (HTTP 401).
type ErrorUnauthorized struct{}

func (e *ErrorUnauthorized) Error() string {
	return "Not authorized"
}

// Venue (symbol) not found (HTTP 404).
type ErrorVenueNotFound struct {
	VenueSymbol string
}

func (e *ErrorVenueNotFound) Error() string {
	return "Venue not found: " + e.VenueSymbol
}

// Stock (symbol) not found in the venue (HTTP 404).
type ErrorStockNotFound struct {
	VenueSymbol string
	StockSymbol string
}

func (e *ErrorStockNotFound) Error() string {
	return fmt.Sprintf("Stock not found: %v (venue: %v)", e.StockSymbol, e.VenueSymbol)
}
