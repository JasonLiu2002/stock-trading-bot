package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/piquette/finance-go/quote"
)

const MOVING_AVG_MULTIPLIER = 1.2

func movingAvgComparison(assetList []string) (buyList []stock, sellList []stock) {
	buyList, sellList = []stock{}, []stock{}

	fmt.Println("Scanning using moving average comparison method...")

	iterator := quote.List(assetList)
	i := 0
	for iterator.Next() {
		quote := iterator.Quote()
		i++
		_, err := alpaca.GetPosition(quote.Symbol)
		if quote.FiftyDayAverage > MOVING_AVG_MULTIPLIER*quote.TwoHundredDayAverage {
			buyList = append(buyList, stock{quote.Symbol, quote.RegularMarketPrice})
		} else if err == nil && quote.FiftyDayAverage < quote.TwoHundredDayAverage {
			sellList = append(sellList, stock{quote.Symbol, quote.RegularMarketPrice})
		}
	}
	fmt.Println("Scan completed")
	return buyList, sellList
}
