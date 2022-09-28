package spotandmargin

//	Необходим для удобного создания подписок
type Cmd struct {
	Time    int64    `json:"time"`
	ID      int64    `json:"id"`
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Payload []string `json:"payload"`
	Auth    authCmd  `json:"auth"`
}
type authCmd struct {
	Method string `json:"method"`
	Key    string `json:"KEY"`
	Sign   string `json:"SIGN"`
}
