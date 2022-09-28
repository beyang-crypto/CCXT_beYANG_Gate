package main

import (
	spotAndMargin "github.com/TestingAccMar/CCXT_beYANG_Gate/gate/spotAndMargin/v4/rest"
)

func main() {
	cfg := &spotAndMargin.Configuration{
		Addr:      spotAndMargin.RestURL,
		ApiKey:    "3ee19c6d0450057413809ea5cea755da",
		SecretKey: "c808a97d1ed945759c9492aa5901660e1e1c1a6d033ca429c7fa3c201938042e",
		DebugMode: true,
	}
	b := spotAndMargin.New(cfg)

	b.GetBalance()
	//a.Auth()

	//a.Start()
	// b.Start()

	// pair := spotAndMargin.GetPair("BTC", "USDT")

	// b.Subscribe(spotAndMargin.ChannelTicker, pair)
	//b.SubscribeOnBalance(spotAndMargin.ChannelBalances)

	// b.On(spotAndMargin.ChannelTicker, handleBestBidPrice)

	//	не дает прекратить работу программы
	forever := make(chan struct{})
	<-forever
}

// func handleBookTicker(symbol string, data spotAndMargin.Tickers) {
// 	log.Printf("Bybit BookTicker  %s: %v", symbol, data)
// }

// func handleBestBidPrice(symbol string, data spotAndMargin.Tickers) {
// 	log.Printf("Bybit BookTicker  %s: BestBidPrice : %s", symbol, data.Result.HighestBid)
// }

// func handleWalletBalanceCoin(data spotAndMargin.WalletBalance) {
// 	for _, coin := range data.Result {
// 		log.Printf("%s :  %s", coin.Currency, coin.Total)
// 	}
// }
