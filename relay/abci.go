package relay

type Abci interface {
	Create(string, string, string) error
	Destroy() error
}

var abciMap map[string]Abci = make(map[string]Abci)

func init() {
	abciMap["eth"] = &EthAbci{}
}
