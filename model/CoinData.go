package model

type CoinData struct {
	CoinData []CoinDatum `json:"data"`
	Status   CMCStatus   `json:"status"`
}

type CMCStatus struct {
	TimeStamp string `json:""`
}

type CoinDatum struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	Slug              string
	Symbol            string `json:"symbol"`
	Logo              string
	DateAdded         string   `json:"date_added"`
	Properties        Property `json:"quote"`
	Tags              []string `json:"tags"`
	MaxSupply         float64  `json:"max_supply"`
	TotalSupply       float64  `json:"total_supply"`
	CirculatingSupply float64  `json:"circulating_supply"`
	Explorers         []string
	Twitters          []string
	Websites          []string
	Chats             []string
	Facebooks         []string
	MessageBoards     []string
	Technicals        []string
	SourceCodes       []string
	Announcements     []string
	BscScan           string
	EthScan           string
	XrpScan           string
	BscContract       string
	EthContract       string
	XrpContract       string
}

type Property struct {
	Dollar MarketValue `json:"USD"`
}
type MarketValue struct {
	LastBlock               BlockchainBlock
	Price                   float64 `json:"price"`
	Volume24                float64 `json:"volume_24h"`
	VolumeChange24          float64 `json:"volume_change_24h"`
	PercentChange24         float64 `json:"percent_change_24h"`
	PercentChange7d         float64 `json:"percent_change_7d"`
	PercentChange30d        float64 `json:"percent_change_30d"`
	PercentChange60d        float64 `json:"percent_change_60d"`
	PercentChange90d        float64 `json:"percent_change_90d"`
	MarketCap               float64 `json:"market_cap"`
	MarketCapDominance      float64 `json:"market_cap_dominance"`
	FullyDilutedMarketPrice float64 `json:"fully_diluted_market_pricevolume_change_1h"`
}
