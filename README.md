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

func main() {
  client := stockfighter.NewClient(*apiKey)

	orderbook, err := client.GetOrderbook(*venueSymbol, *stockSymbol)
	if err != nil {
		panic(err)
	}
	
	// ... 

	quote, err := client.GetQuote(*venueSymbol, *stockSymbol)
	if err != nil {
		panic(err)
	}

	// ...
}
```

## References

See [GoDoc](https://godoc.org/gpk.io/stockfighter).
