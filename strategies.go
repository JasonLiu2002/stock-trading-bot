package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/piquette/finance-go/quote"
)

const MOVING_AVG_THRESHOLD = 1.2

func movingAvgComparison(assetList []alpaca.Asset) (buyList []stock, sellList []stock) {
	buyList, sellList = []stock{}, []stock{}

	fmt.Println("Algo trading began using moving average comparison method...")
	for _, asset := range assetList {
		_, err := alpaca.GetPosition(asset.Symbol)
		quote, _ := quote.Get(asset.Symbol)
		if quote != nil {
			if quote.FiftyDayAverage > MOVING_AVG_THRESHOLD*quote.TwoHundredDayAverage {
				fmt.Printf("Identified %v stock to buy. \n", asset.Symbol)
				buyList = append(buyList, stock{asset.Symbol, quote.RegularMarketPrice})
			} else if err == nil && quote.FiftyDayAverage < quote.TwoHundredDayAverage {
				fmt.Printf("Identified %v stock to sell. \n", asset.Symbol)
				sellList = append(sellList, stock{asset.Symbol, quote.RegularMarketPrice})
			}
		}
	}
	return buyList, sellList
}
