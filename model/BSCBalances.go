package model

type BSCBalances struct {
	Balances []BSCBalance
}

type BSCBalance struct {
	Address    string
	Contract   string
	Amount     float64
	USDConvert float64
	CoinDatum  CoinDatum
}
