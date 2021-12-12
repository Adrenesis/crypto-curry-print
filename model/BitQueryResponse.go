package model

type BitQueryBSCBalanceResult struct {
	Data BitQueryBSCDataBank `json:"data"`
}

type BitQueryBSCDataBank struct {
	//s string `json:"address"`
	Properties BitQueryBSCDataApi `json:"ethereum"`
}

type BitQueryBSCDataApi struct {
	//s string `json:"address"`
	API []BitQueryBSCProperties `json:"address"`
}

type BitQueryBSCProperties struct {
	Balances []BitQueryBSCBalances `json:"balances"`
}

type BitQueryBSCBalances struct {
	Token BitQueryBSCToken `json:"currency"`
	//Token string
	Amount float64 `json:"value"`
}

type BitQueryBSCToken struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}

type BitQueryBSCQuoteResult struct {
	Data BitQueryBSCQuoteDataBank `json:"data"`
}

type BitQueryBSCQuoteDataBank struct {
	//s string `json:"address"`
	Properties BitQueryBSCQuoteDataApi `json:"ethereum"`
}

type BitQueryBSCQuoteDataApi struct {
	//s string `json:"address"`
	API []BitQueryBSCQuoteTrades `json:"DexTrades"`
}

type BitQueryBSCQuoteTrades struct {
	Currency      BitQueryBSCToken `json:"baseCurrency"`
	QuoteCurrency BitQueryBSCToken `json:"quoteCurrency"`
	Price         float64          `json:"quotePrice"`
	Block         BlockchainBlock  `json:"block"`
}

type BlockchainBlock struct {
	Height    int64             `json:"height"`
	TimeStamp BitQueryTimeStamp `json:"timestamp"`
	Network   string
}

type BitQueryTimeStamp struct {
	ISO8601 string `json:"iso8601"`
}
