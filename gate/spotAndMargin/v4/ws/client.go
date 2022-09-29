package spotandmargin

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"      //  Для вытаскивания одного значения из файла json
	"github.com/chuckpreslar/emission" // Эмитер необходим для удобного выполнения функции в какой-то момент
	"github.com/goccy/go-json"         // для создания собственных json файлов и преобразования json в структуру
	"github.com/gorilla/websocket"
)

const (
	HostWebsocketURL = "wss://api.gateio.ws/ws/v4/"
)

const (
	//https://www.gate.io/docs/developers/apiv4/ws/en/#tickers-channel
	ChannelTicker = "spot.tickers"

	// https://www.gate.io/docs/developers/apiv4/ws/en/#spot-balance-channel
	ChannelBalances = "spot.balances"
)

type Configuration struct {
	Addr      string `json:"addr"`
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	DebugMode bool   `json:"debug_mode"`
}

type GateWS struct {
	cfg  *Configuration
	conn *websocket.Conn

	mu            sync.RWMutex
	subscribeCmds []Cmd //	сохраняем все подписки у данной биржи, чтоб при переподключении можно было к ним повторно подключиться

	emitter *emission.Emitter
}

func (b *GateWS)  GetPair(coin1 string, coin2 string) string {
	return strings.ToUpper(coin1 + "_" + coin2)
}

func New(config *Configuration) *GateWS {

	// 	потом тут добавятся различные другие настройки
	b := &GateWS{
		cfg:     config,
		emitter: emission.NewEmitter(),
	}
	return b
}

func (b *GateWS) Subscribe(args ...string) {
	switch len(args) {
	case 1:
		b.Subscribe1(args[0])
	case 2:
		b.Subscribe2(args[0], args[1])
	default:
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
				"File": "client.go",
				"Functions" : "(b *GateWS) Subscribe(args ...string
				"Exchange" : "Gate",
				"Data" : [%v],
				"Comment" : "Слишком много аргументов"
			}`, args)
		log.Fatal()
	}
}

func (b *GateWS) Subscribe1(channel string) {
	ts := time.Now().Unix()*1000 + 10000
	auth := authCmd{
		Method: "api_key",
		Key:    b.cfg.ApiKey,
		Sign:   b.sign(channel, "subscribe", ts),
	}
	cmd := Cmd{
		Time:    ts,
		Channel: channel,
		Event:   "subscribe",
		Auth:    auth,
	}
	b.subscribeCmds = append(b.subscribeCmds, cmd)
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 1")
	}
	b.SendCmd(cmd)
}

func (b *GateWS) Subscribe2(channel string, coin ...string) {
	ts := time.Now().Unix()*1000 + 10000
	auth := authCmd{
		Method: "api_key",
		Key:    b.cfg.ApiKey,
		Sign:   b.sign(channel, "subscribe", ts),
	}
	cmd := Cmd{
		Time:    ts,
		Channel: channel,
		Event:   "subscribe",
		Payload: coin,
		Auth:    auth,
	}
	b.subscribeCmds = append(b.subscribeCmds, cmd)
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 1")
	}
	b.SendCmd(cmd)
}

func (b *GateWS) sign(channel, event string, t int64) string {
	message := fmt.Sprintf("channel=%s&event=%s&time=%d", channel, event, t)
	h2 := hmac.New(sha512.New, []byte(b.cfg.SecretKey))
	io.WriteString(h2, message)
	return hex.EncodeToString(h2.Sum(nil))
}

//	отправка команды на сервер в отдельной функции для того, чтобы при переподключении быстро подписаться на все предыдущие каналы
func (b *GateWS) SendCmd(cmd Cmd) {
	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
				"File": "client.go",
				"Functions" : "(b *GateWS) sendCmd(cmd Cmd)",
				"Function where err" : "json.Marshal",
				"Exchange" : "Gate",
				"Data" : [%v],
				"Error" : %s
			}`, cmd, err)
		log.Fatal()
	}
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 2")
	}
	b.Send(string(data))
}

func (b *GateWS) Send(msg string) (err error) {
	defer func() {
		// recover необходим для корректной обработки паники
		if r := recover(); r != nil {
			if err != nil {
				log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
						"File": "client.go",
						"Functions" : " (b *GateWS) Send(msg string) (err error)",
						"Function where err" : "b.conn.WriteMessage",
						"Exchange" : "Gate",
						"Data" : [websocket.TextMessage, %s],
						"Error" : %s,
						"Recover" : %v
					}`, msg, err, r)
				log.Fatal()
			}
			err = errors.New(fmt.Sprintf("GateWs send error: %v", r))
		}
	}()
	if b.cfg.DebugMode {
		log.Printf("Отправка сообщения на сервер. текст сообщения:%s", msg)
	}

	err = b.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return
}

// подключение к серверу и постоянное чтение приходящих ответов
func (b *GateWS) Start() error {
	if b.cfg.DebugMode {
		log.Printf("Начало подключения к серверу")
	}
	b.connect()

	cancel := make(chan struct{})

	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				b.ping()
			case <-cancel:
				return
			}
		}
	}()

	go func() {
		defer close(cancel)

		for {
			_, data, err := b.conn.ReadMessage()
			if err != nil {

				if websocket.IsCloseError(err, 1006) {
					b.closeAndReconnect()
					//Необходим вызв SubscribeToTicker в отдельной горутине, рекурсия, думаю, тут неуместна
					log.Printf("Status: INFO	ошибка 1006 начинается переподключение к серверу")

				} else {
					b.conn.Close()
					log.Printf(`
						{
							"Status" : "Error",
							"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
							"File": "client.go",
							"Functions" : "(b *GateWS) Start() error",
							"Function where err" : "b.conn.ReadMessage",
							"Exchange" : "Gate",
							"Error" : %s
						}`, err)
					log.Fatal()
				}

			} else {
				b.messageHandler(data)
			}

		}
	}()

	return nil
}

func (b *GateWS) connect() {

	c, _, err := websocket.DefaultDialer.Dial(b.cfg.Addr, nil)
	if err != nil {
		log.Printf(`{
						"Status" : "Error",
						"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
						"File": "client.go",
						"Functions" : "(b *GateWS) connect()",
						"Function where err" : "websocket.DefaultDialer.Dial",
						"Exchange" : "Gate",
						"Data" : [%s, nil],
						"Error" : %s
					}`, b.cfg.Addr, err)
		log.Fatal()
	}
	b.conn = c
	for _, cmd := range b.subscribeCmds {
		b.SendCmd(cmd)
	}
}

func (b *GateWS) closeAndReconnect() {
	b.conn.Close()
	b.connect()
}

func (b *GateWS) ping() {
	ts := time.Now().Unix()*1000 + 10000
	defer func() {
		if r := recover(); r != nil {
			log.Printf("GateWs ping error: %v", r)
		}
	}()

	ping := fmt.Sprintf(`
	{
		"time": %d,
		"channel" : "spot.ping"
	}`, ts)
	err := b.conn.WriteMessage(websocket.TextMessage, []byte(ping))
	if err != nil {
		log.Printf("GateWs ping error: %v", err)
	}
}

func (b *GateWS) messageHandler(data []byte) {
	if b.cfg.DebugMode {
		log.Printf("GateWs %v", string(data))
	}

	channel, _ := jsonparser.GetString(data, "channel")
	channelArr := strings.Split(channel, ".")
	switch channelArr[1] {
	case "pong":
		//	ничего не надо
	case "tickers":
		var ticker Tickers
		err := json.Unmarshal(data, &ticker)
		if err != nil {
			log.Printf(`
				{
					"Status" : "Error",
					"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
					"File": "client.go",
					"Functions" : "(b *GateWS) messageHandler(data []byte)",
					"Function where err" : "json.Unmarshal",
					"Exchange" : "Gate",
					"Comment" : %s to BookTicker struct,
					"Error" : %s
				}`, string(data), err)
			log.Fatal()
		}
		b.processTicker(channelArr[1], ticker)
	case "balances":
		event, _ := jsonparser.GetString(data, "event")
		if event != "subscribe" {
			var walletBalance WalletBalance
			err := json.Unmarshal(data, &walletBalance)
			if err != nil {
				log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
						"File": "client.go",
						"Functions" : "(b *GateWS) messageHandler(data []byte)",
						"Function where err" : "json.Unmarshal",
						"Exchange" : "Gate",
						"Comment" : %s to BookTicker struct,
						"Error" : %s
					}`, string(data), err)
				log.Fatal()
			}
			b.processWalletBalance(channelArr[1], walletBalance)
		}
	default:
		log.Printf(`
			{
				"Status" : "INFO",
				"Path to file" : "CCXT_beYANG_Gate/gate/spotAndMargin/v4/ws",
				"File": "client.go",
				"Functions" : "(b *GateWS) messageHandler(data []byte)",
				"Exchange" : "Gate",
				"Comment" : "Ответ от неизвестного канала"
				"Message" : %s
			}`, string(data))
		log.Fatal()
	}
}
