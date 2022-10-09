package rest

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Configuration struct {
	Addr      string `json:"addr"`
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	DebugMode bool   `json:"debug_mode"`
}

type GateRest struct {
	cfg *Configuration
}

const (
	RestURL = "https://api.gateio.ws"
)

const (
	WalletDepositsrecords = "/api/v4/wallet/deposits"
)

func New(config *Configuration) *GateRest {

	// 	потом тут добавятся различные другие настройки
	b := &GateRest{
		cfg: config,
	}
	return b
}

func (ex *GateRest) GetBalance() interface{} {

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	//	https://www.gate.io/docs/developers/apiv4/#retrieve-withdrawal-records
	//	получение времяниz
	ts := time.Now().Unix()
	apiKey := ex.cfg.ApiKey
	secretKey := ex.cfg.SecretKey
	hasher := sha512.New()
	hasher.Write([]byte(""))
	hashed_payload := hex.EncodeToString(hasher.Sum(nil)) // да, так и должно быть, что пусто
	url := WalletDepositsrecords
	parms := fmt.Sprintf("%s\n%s\n%s\n%s\n%d", "GET", url, "", hashed_payload, ts)
	log.Printf(parms)
	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write([]byte(parms))
	sign := hex.EncodeToString(mac.Sum(nil))
	log.Printf("===\n===" + sign)
	url = ex.cfg.Addr + url
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", apiKey)
	req.Header.Set("Timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("SIGN", sign)
	//	код для вывода полученных данных
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Do(req)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if ex.cfg.DebugMode {
		log.Printf("STATUS: DEBUG\tEXCHANGE: Gate\tAPI: Rest\tGate WalletDeposits %v", string(data))
	}

	var walletBalance WalletDeposits
	err = json.Unmarshal(data, &walletBalance)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_Gate/gate/spotAndMargin/v4/rest",
				"File": "client.go",
				"Functions" : "ex *GateRest) GetBalance() interface{}",
				"Function where err" : "json.Unmarshal",
				"Exchange" : "Gate",
				"Comment" : %s to WalletDeposits struct,
				"Error" : %s
			}`, string(data), err)
		log.Fatal()
	}

	return walletBalance

}
