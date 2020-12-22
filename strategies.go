package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/piquette/finance-go/quote"
)

const MOVING_AVG_MULTIPLIER = 1.2
const VWAP_MULTIPLIER = 1.05

func movingAvgComparison(assetList []string) (buyList []stock, sellList []stock) {
	buyList, sellList = []stock{}, []stock{}

	fmt.Println("Scanning using moving average comparison method...")

	iterator := quote.List(assetList)
	for iterator.Next() {
		quote := iterator.Quote()
		if quote.IsTradeable {
			_, err := alpaca.GetPosition(quote.Symbol)
			if quote.FiftyDayAverage > MOVING_AVG_MULTIPLIER*quote.TwoHundredDayAverage {
				buyList = append(buyList, stock{quote.Symbol, quote.Ask})
			} else if err == nil && quote.FiftyDayAverage < quote.TwoHundredDayAverage {
				sellList = append(sellList, stock{quote.Symbol, quote.Bid})
			}
		}
	}

	fmt.Println("Scan completed")
	return buyList, sellList
}

// Sell when above the VWAP, buy when below the VWAP
func volumeWeightedAveragePrice(assetList []string) (buyList []stock, sellList []stock) {
	buyList, sellList = []stock{}, []stock{}

	fmt.Println("Scanning using volume weighted average price method...")

	dict := calcVwap(assetList)
	iterator := quote.List(assetList)
	for iterator.Next() {
		quote := iterator.Quote()
		if quote.IsTradeable {
			_, err := alpaca.GetPosition(quote.Symbol)
			if quote.Ask < (1.0/VWAP_MULTIPLIER)*dict[quote.Symbol] { // Buy condition
				fmt.Printf("Buying: %v. VWAP: %v, Ask: %v \n", quote.Symbol, dict[quote.Symbol], quote.Ask)
				buyList = append(buyList, stock{quote.Symbol, quote.Ask})
			} else if err == nil && quote.Bid > VWAP_MULTIPLIER*dict[quote.Symbol] { // Sell condition
				fmt.Printf("Selling: %v. VWAP: %v, Bid: %v \n", quote.Symbol, dict[quote.Symbol], quote.Bid)
				sellList = append(sellList, stock{quote.Symbol, quote.Bid})
			}
		}
	}

	fmt.Println("Scan completed")
	return buyList, sellList
}
