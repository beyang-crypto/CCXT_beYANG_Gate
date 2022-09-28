package spotandmargin

func (b *GateWS) processTicker(symbol string, data Tickers) {
	b.Emit(ChannelTicker, symbol, data)
}

func (b *GateWS) processWalletBalance(symbol string, data WalletBalance) {
	b.Emit(ChannelBalances, data)
}
