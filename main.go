package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/shopspring/decimal"
)

type stock struct {
	symbol string
	price  float64
}

func init() {
	alpaca.SetBaseUrl("https://paper-api.alpaca.markets")
}

func main() {
	assetList := getAssets()
	clock, _ := alpaca.GetClock()
	account, err := alpaca.GetAccount()
	if err != nil {
		panic(err)
	}
	buyingPower := account.BuyingPower

	for true {
		if clock.IsOpen {
			buyList, sellList := movingAvgComparison(assetList)

			manageStockPurchases(buyList, buyingPower)
			manageStockSales(sellList)
		}
	}
}

func manageStockPurchases(buyList []stock, buyingPower decimal.Decimal) {
	if buyingPower.GreaterThan(decimal.NewFromInt(0)) {
		moneyPerSymbol := buyingPower.Div(decimal.NewFromInt(int64(len(buyList))))

		for _, asset := range buyList {
			numShares := moneyPerSymbol.Div(decimal.NewFromFloat(asset.price).RoundDown(0))
			if numShares.GreaterThan(decimal.NewFromInt(0)) {
				fmt.Printf("Buying %v shares of %v at %v each. \n", numShares, asset.symbol, asset.price)
				buy(asset.symbol, numShares)
			}
		}
	}
}

func manageStockSales(sellList []stock) {
	for _, asset := range sellList {
		position, err := alpaca.GetPosition(asset.symbol)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Selling %v shares of %v at %v each. \n", position.Qty, asset.symbol, asset.price)
		sell(asset.symbol, position.Qty)
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
		if asset.Tradable {
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

func buy(symbol string, numShares decimal.Decimal) {
	alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Qty:         numShares,
		Side:        alpaca.Buy,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
	})
}
