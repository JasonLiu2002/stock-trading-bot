# **Stock Trading Bot**
### About
An algorithmic stock trading bot written in Golang that uses the [Alpaca API](https://alpaca.markets/)

### Algorithms
- Volume-weighted average price
- 50-day and 200-day moving average comparisons

### Setup
1. Register for an Alpaca account [here](https://app.alpaca.markets/signup)\
2. Navigate to paper trading and generate your API keys\
3. Set your private and public keys as environment variables -- I do this through a .zshrc or .bash_profile:\
`export APCA_API_KEY_ID="XXXXXXXXXXXXXXXXXXXX"`\
`export APCA_API_SECRET_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`
4. Install the necessary dependencies\
   `go get github.com/alpacahq/alpaca-trade-api-go/alpaca`\
   `go get github.com/shopspring/decimal`\
   `go get github.com/piquette/finance-go/quote`
5. `go run .`
6. Enjoy :)

### To-do
- Add sentiment analysis
- Combine VWAP with TWAP
- Add mean reversion
- Improve web scraping method for stocks


