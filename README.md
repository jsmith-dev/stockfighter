# Stockfighter

[![GoDoc](https://godoc.org/gpk.io/stockfighter?status.svg)](https://godoc.org/gpk.io/stockfighter)

## Install

```bash
go get gpk.io/stockfighter
```

## Example

```go
package main

import "gpk.io/stockfighter.v0"

const (
	apiKey = "your_stockfighter_api_key"
	venue = "ABCD"
	stock = "XYZ"
)

func main() {
	client := stockfighter.NewClient(apiKey)

	orderbook, err := client.GetOrderbook(venue, stock)
	if err != nil {
		panic(err)
	}
	
	// ... 
	
	quote, err := client.GetQuote(venue, stock)
	if err != nil {
		panic(err)
	}
	
	// ...
}
```

## Tests

To run unit tests, run:

```bash
API_KEY=your_stockfighter_api_key go test
```

## References

See [GoDoc](https://godoc.org/gpk.io/stockfighter).
