package model

type ApiKeyData struct {
	CMCApiKey      string `json:"coin_market_cap_api_key"`
	BSCApiKey      string `json:"bscscan_api_key"`
	BitQueryApiKey string `json:"bitquery_api_key"`
}
