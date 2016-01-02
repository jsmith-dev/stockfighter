/*
Package stockfighter is a wrapper for Stockfighter API (v1.0).

    var apiKey = "your_stockfighter_api_key"
    var venue = "ABCD"
    var stock = "XYZ"

    client := stockfighter.NewClient(apiKey)

    err := client.Ping()
    if err != nil {
        panic(err)
    }

    orderbook, err := client.GetStockOrderbook(venue, stock)
    if err != nil {
        panic(err)
    }

    // ...

For more information on Stockfighter API, see https://starfighter.readme.io/v1.0/docs.
*/
package stockfighter
