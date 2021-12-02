package model

type ApiKeyData struct {
	ApiKey string `json:"api_key"`
}

type CoinData struct {
	CoinData []CoinDatum `json:"data"`
}

type CoinDatum struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	Symbol     string   `json:"symbol"`
	DateAdded  string   `json:"date_added"`
	Properties Property `json:"quote"`
}

type Property struct {
	Dollar MarketValue `json:"USD"`
}
type MarketValue struct {
	Price    float64 `json:"price"`
	Volume24 float64 `json:"volume_24h"`
}
