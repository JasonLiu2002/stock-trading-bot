package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
)

func calcVwap(assetList []string) map[string]float64 {
	limit := 500
	dict := make(map[string]float64)
	for _, asset := range assetList {
		bars, err := alpaca.GetSymbolBars(asset, alpaca.ListBarParams{
			Timeframe: "minute",
			Limit:     &limit,
		})
		if err != nil {
			panic(err)
		}
		sumPV := 0.0
		sumVolume := 0
		for _, bar := range bars {
			sumPV += float64(bar.Volume) * float64(bar.High+bar.Low+bar.Close) / 3.0
			sumVolume += int(bar.Volume)
		}
		dict[asset] = sumPV / float64(sumVolume)
	}
	return dict
}
