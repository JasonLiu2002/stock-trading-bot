package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/shopspring/decimal"
	"time"
)

type stock struct {
	symbol string
	price  float64
}

func init() {
	alpaca.SetBaseUrl("https://paper-api.alpaca.markets")
}

func main() {
	// TODO: Scrape these assets from websites instead
	//assetList := getAssets()
	assetList := []string{
		"NIO", "AAPL", "AAL", "BAC", "GE", "PLTR", "FEYE", "F", "TSLA", "CCL", "INTC", "AMD", "WFC", "ITUB", "LAZR",
		"T", "FCEL", "PLUG", "PFE", "FUBO", "C", "MSFT", "OXY", "NOK", "BBD", "JPM", "XOM", "NCLH", "MRNA", "VALE",
		"UBER", "SWN", "QS", "UAL", "LI", "ET", "M", "ABEV", "BP", "MRO", "XPEV", "BA", "HBAN", "BB", "KGC", "KMI",
		"BFT", "AJRD", "SPCE", "RYCEY", "SWI", "DAL", "ABNB", "PBR", "GOLD", "SAN", "JMIA", "LKNCY", "SIRI", "CSCO",
		"IQ", "RP", "FB", "NKE", "BIDU", "ZNGA", "SDC", "COTY", "MGNI", "ORCL", "LYG", "VZ", "GM", "MS", "AZN", "WKHS",
		"AI", "DKNG", "X", "MU", "AUY", "PTON", "GLIBA", "EDIT", "ATUS", "PINS", "HL", "APHA", "VER", "LUMN", "KO",
		"SRNE", "FCX", "CVX", "DIS", "BABA", "SQ", "ZM", "CMCSA", "COP", "COP", "NLY", "ZM", "VTRS", "SPWR", "WORK",
		"WISH", "CLF", "TWTR", "SU", "GILD", "HPE", "VLDR", "NKLA", "KEY", "MO", "HPQ", "CRM", "JBLU", "DVN", "RDS-A",
		"SABR", "MRK", "HAL", "RF", "SNAP", "HST", "MGM", "SAVE", "CRWD", "BMY", "PACB", "AG", "BSX", "BILI", "EPD",
		"CHWY", "GS", "FSR", "MRVL", "LYFT", "APA", "WMT", "WMB", "KR", "VOD",
		"BRPA", "BRPAU", "SNGX", "OCGN", "CCRC", "SOS", "NH", "JG", "CYCC", "YGMZ", "AQMS", "OLB", "AMTX", "MVIS",
		"SGTX", "CLA.U", "VISL", "RIOT", "ELYS", "CLA", "IPV.U", "FRSX", "NK", "VERO", "WNW", "GRAY", "PTN", "STIC",
		"UAVS", "GNPX", "KOPN", "JAGX", "QS", "CMPI", "AHT", "EYES", "QK", "AMPE", "IPV", "LPTH", "NRIX", "SPWH", "ONVO",
		"GME", "HPK", "MFH", "HUSN", "SQFT", "VSTO", "ONTX", "VUZI", "TIGR", "IHT", "LGVW.U", "SURF", "BEAM", "ONDS",
		"CLNE", "VLDR", "TLSA", "BLNK", "CGROU", "ADXN", "LGVW", "SOL", "EAR", "NISN", "DMTK", "HARP", "CBAT", "STPK.U", "CNSP",
		"LPTX", "ALVR", "USIO", "SGMO", "EQOS", "NES", "OCN", "FCEL", "INDO", "BCRX", "SYRS", "RGLS", "TCS", "SPRO", "SLGG", "ITP",
		"EH", "AQB", "LOAK.U", "NXTD", "LAC", "WBAI", "FUV", "MARA", "SEEL", "MTLS", "CCCC", "TCON",
	}
	clock, _ := alpaca.GetClock()

	var account *alpaca.Account
	var err error
	var todaysDate time.Time

	updateTime := time.Now()
	firstTime := true
	if !clock.IsOpen {
		fmt.Println("Markets are closed, waiting...")
	}

	for true {
		if clock.IsOpen {
			if firstTime {
				fmt.Println("Markets are open!")
				year, month, day := time.Now().Date()
				location, _ := time.LoadLocation("GMT")
				todaysDate = time.Date(year, month, day, 0, 0, 0, 0, location)
			}
			// Scan once a minute
			if time.Now().Sub(updateTime) > time.Duration(60e9) || firstTime {
				firstTime = false
				updateTime = time.Now()
				account, err = alpaca.GetAccount()
				if err != nil {
					panic(err)
				}
				fmt.Println("Updated account info")
				fmt.Printf("Equity: %v \n", account.Equity)
				fmt.Printf("Buying power: %v \n", account.BuyingPower)

				buyList, sellList := volumeWeightedAveragePrice(assetList, todaysDate)
				//buyList, sellList := movingAvgComparison(assetList)

				manageStockPurchases(buyList, account.BuyingPower.Div(decimal.NewFromInt(4)))
				manageStockSales(sellList)
				fmt.Println()
			}
		}
	}
}

func manageStockPurchases(buyList []stock, buyingPower decimal.Decimal) {
	if buyingPower.GreaterThan(decimal.NewFromInt(0)) && len(buyList) > 0 {
		moneyPerSymbol := buyingPower.Div(decimal.NewFromInt(int64(len(buyList))))

		for _, asset := range buyList {
			numShares := moneyPerSymbol.Div(decimal.NewFromFloat(asset.price)).RoundDown(0)
			if numShares.GreaterThan(decimal.NewFromInt(0)) {
				fmt.Printf("Buying %v shares of %v at %v each. \n", numShares, asset.symbol, asset.price)
				_, err := alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
					AssetKey:    &asset.symbol,
					Qty:         numShares,
					Side:        alpaca.Buy,
					Type:        alpaca.Market,
					TimeInForce: alpaca.Day,
				})
				if err != nil {
					panic(err)
				}
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
		_, sellErr := alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
			AssetKey:    &asset.symbol,
			Qty:         position.Qty,
			Side:        alpaca.Sell,
			Type:        alpaca.Market,
			TimeInForce: alpaca.Day,
		})
		if sellErr != nil {
			panic(sellErr)
		}
	}
}

func getAssets() []string {
	status := "active"
	assets, err := alpaca.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	var tradeableAssets []string
	for _, asset := range assets {
		if asset.Tradable {
			tradeableAssets = append(tradeableAssets, asset.Symbol)
		}
	}

	return tradeableAssets
}
