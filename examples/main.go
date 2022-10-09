package main

import (
	"log"

	spotAndMarginWs "github.com/TestingAccMar/CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws"
)

func main() {
	cfg := &spotAndMarginWs.Configuration{
		Addr:      spotAndMarginWs.HostWebsocketURL,
		ApiKey:    "",
		SecretKey: "",
		DebugMode: true,
	}
	b := spotAndMarginWs.New(cfg)

	b.Start()

	pair1 := b.GetPair("BTC", "USDT")
	pair2 := b.GetPair("eth", "USDT")

	b.Subscribe(spotAndMarginWs.ChannelTicker, []string{pair1, pair2})

	b.On(spotAndMarginWs.ChannelTicker, handleBestBidPrice)

	//	не дает прекратить работу программы
	forever := make(chan struct{})
	<-forever
}

// func handleBookTicker(symbol string, data spotAndMargin.Tickers) {
// 	log.Printf("Bybit BookTicker  %s: %v", symbol, data)
// }

func handleBestBidPrice(name string, symbol string, data spotAndMarginWs.Tickers) {
	log.Printf("%s BookTicker  %s: BestBidPrice : %s", name, symbol, data.Result.HighestBid)
}

// func handleWalletBalanceCoin(data spotAndMargin.WalletBalance) {
// 	for _, coin := range data.Result {
// 		log.Printf("%s :  %s", coin.Currency, coin.Total)
// 	}
// }
