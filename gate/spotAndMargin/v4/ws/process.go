package spotandmargin

func (b *GateWS) processTicker(name string, symbol string, data Tickers) {
	b.Emit(ChannelTicker, name, symbol, data)
}

func (b *GateWS) processWalletBalance(name string, symbol string, data WalletBalance) {
	b.Emit(ChannelBalances, name, data)
}
