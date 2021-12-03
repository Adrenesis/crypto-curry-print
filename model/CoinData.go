package model

type ApiKeyData struct {
	ApiKey string `json:"api_key"`
}

type CoinData struct {
	CoinData []CoinDatum `json:"data"`
}

type CoinDatum struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Slug          string
	Symbol        string `json:"symbol"`
	Logo          string
	DateAdded     string   `json:"date_added"`
	Properties    Property `json:"quote"`
	Explorers     []string
	Twitters      []string
	Websites      []string
	Chats         []string
	Facebooks     []string
	MessageBoards []string
	Technicals    []string
	SourceCodes   []string
	Announcements []string
	BscScan       string
	EthScan       string
	XrpScan       string
	BscContract   string
	EthContract   string
	XrpContract   string
}

type Property struct {
	Dollar MarketValue `json:"USD"`
}
type MarketValue struct {
	Price    float64 `json:"price"`
	Volume24 float64 `json:"volume_24h"`
}
