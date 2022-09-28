package spotandmargin

// https://www.gate.io/docs/developers/apiv4/#retrieve-ticker-information
type Tickers struct {
	Time    int           `json:"time"`
	Channel string        `json:"channel"`
	Event   string        `json:"event"`
	Result  resultTickers `json:"result"`
}
type resultTickers struct {
	CurrencyPair     string `json:"currency_pair"`
	Last             string `json:"last"`
	LowestAsk        string `json:"lowest_ask"`
	HighestBid       string `json:"highest_bid"`
	ChangePercentage string `json:"change_percentage"`
	BaseVolume       string `json:"base_volume"`
	QuoteVolume      string `json:"quote_volume"`
	High24H          string `json:"high_24h"`
	Low24H           string `json:"low_24h"`
}

//https://www.gate.io/docs/developers/apiv4/ws/en/#client-subscription-9
type WalletBalance struct {
	Time    int                   `json:"time"`
	Channel string                `json:"channel"`
	Event   string                `json:"event"`
	Result  []resultWalletBalance `json:"result"`
}
type resultWalletBalance struct {
	Timestamp   string `json:"timestamp"`
	TimestampMs string `json:"timestamp_ms"`
	User        string `json:"user"`
	Currency    string `json:"currency"`
	Change      string `json:"change"`
	Total       string `json:"total"`
	Available   string `json:"available"`
}
