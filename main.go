package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/piquette/finance-go/quote"
	"github.com/shopspring/decimal"
	"strings"
)

var clock alpaca.Clock
var account alpaca.Account

func main() {
	alpaca.SetBaseUrl("https://paper-api.alpaca.markets")

	clock, _ := alpaca.GetClock()
	assetList := getAssets()

	for true {
		if clock.IsOpen {
			fmt.Println("Scan began...")
			for _, asset := range assetList {
				position, err := alpaca.GetPosition(asset.Symbol)
				quote, _ := quote.Get(asset.Symbol)
				if quote != nil {
					if err != nil {
						if quote.FiftyDayAverage > 1.2*quote.TwoHundredDayAverage {
							//buy(asset.Symbol, 5)
							fmt.Println("Buying 5 shares of " + asset.Symbol)
						}
					} else {
						if quote.FiftyDayAverage < quote.TwoHundredDayAverage {
							//sell(asset.Symbol, position.Qty)
							fmt.Println("Selling " + position.Qty.String() + " shares of " + asset.Symbol)
						}
					}
				}
			}
		}
	}
}

func getAssets() []alpaca.Asset {
	status := "active"
	assets, err := alpaca.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	tradableAssets := []alpaca.Asset{}
	for _, asset := range assets {
		if asset.Tradable && !strings.Contains(asset.Symbol, "-") && !strings.Contains(asset.Symbol, ".") {
			tradableAssets = append(tradableAssets, asset)
		}
	}

	return tradableAssets
}

func sell(symbol string, numShares decimal.Decimal) {
	alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Qty:         numShares,
		Side:        alpaca.Sell,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
	})
}

func buy(symbol string, numShares int64) {
	alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Qty:         decimal.NewFromInt(numShares),
		Side:        alpaca.Buy,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
	})
}
