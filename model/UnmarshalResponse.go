package model

type UnmarshalBSCBalanceResult []UnmarshalBSCBalanceData

type UnmarshalBSCBalanceData struct {
	//s string `json:"address"`
	Name     string  `json:"contract_name"`
	Symbol   string  `json:"contract_ticker_symbol"`
	Decimals float64 `json:"contract_decimals"`
	Contract string  `json:"contract_address"`
	Amount   string  `json:"balance"`
}
