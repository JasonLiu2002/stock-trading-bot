package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

const MOVING_AVG_MULTIPLIER = 1.2

func movingAvgComparison(assetList []string) (buyList []stock, sellList []stock) {
	buyList, sellList = []stock{}, []stock{}

	fmt.Println("Algo trading began using moving average comparison method...")

	var quoteList []finance.Quote

	/**
	 * TODO: Find some way to make this not sketchy
	 * The quote.List method does not seem to be able to query more than 1500 symbols at a time
	 * One solution is to be more selective with which stocks we will query (which will also save time)
	 */
	for i := 0; i < 9; i++ {
		var smallAssetList []string
		for j := i * 1000; j < (i+1)*1000; j++ {
			smallAssetList = append(smallAssetList, assetList[j])
		}
		iterator := quote.List(smallAssetList)
		for iterator.Next() {
			quoteList = append(quoteList, *iterator.Quote())
		}
	}

	fmt.Println(len(quoteList))

	for _, quote := range quoteList {
		_, err := alpaca.GetPosition(quote.Symbol)
		if quote.FiftyDayAverage > MOVING_AVG_MULTIPLIER*quote.TwoHundredDayAverage {
			fmt.Printf("Identified %v stock to buy. \n", quote.Symbol)
			buyList = append(buyList, stock{quote.Symbol, quote.RegularMarketPrice})
		} else if err == nil && quote.FiftyDayAverage < quote.TwoHundredDayAverage {
			fmt.Printf("Identified %v stock to sell. \n", quote.Symbol)
			sellList = append(sellList, stock{quote.Symbol, quote.RegularMarketPrice})
		}
	}

	return buyList, sellList
}
