package rest

import "log"

//https://www.gate.io/docs/developers/apiv4/#retrieve-deposit-records
type WalletDeposits []struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Currency  string `json:"currency"`
	Address   string `json:"address"`
	Txid      string `json:"txid"`
	Amount    string `json:"amount"`
	Memo      string `json:"memo"`
	Status    string `json:"status"`
	Chain     string `json:"chain"`
}

func GateToWalletBalance(data interface{}) WalletDeposits {
	bt, ok := data.(WalletDeposits)
	if !ok {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/rest",
				"File": "response.go",
				"Functions" : "GateToWalletBalance(data interface{}) WalletDeposits
				"Exchange" : "Gate",
				"Comment" : "Ошибка преобразования %v в WalletDeposits"
			}`, data)
		log.Fatal()
	}
	return bt
}
