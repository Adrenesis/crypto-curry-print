package model

import "strings"

func GetCoinDataByBscContracts(coinData CoinData) map[string]CoinDatum {
	//coinDataIndex := make(map[string]int64)
	coinDataByContract := make(map[string]CoinDatum)
	for i := 0; i < len(coinData.CoinData); i++ {
		//coinDataIndex[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i].Id
		coinDataByContract[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i]
	}
	return coinDataByContract
}
