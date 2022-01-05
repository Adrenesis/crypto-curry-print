package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func ReadBitQueryPriceResultJson(contract string, filename string) CoinData {
	// Open our jsonFile
	jsonFile, err := os.Open("BSCPriceTemp.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(ConvertToISO8601(time.Now()), err)
	}
	fmt.Println(ConvertToISO8601(time.Now()), "Successfully Opened ", "BSCPriceTemp.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var bitQueryResult BitQueryBSCQuoteResult

	errJson := json.Unmarshal(byteValue, &bitQueryResult)
	if errJson != nil {
		fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
		log.Fatal(errJson)
	}

	if len(bitQueryResult.Data.Properties.API) == 0 || bitQueryResult.Data.Properties.API == nil {
		fmt.Sprintf("Failed to request prices, opening backup...")
		jsonFile.Close()

		jsonFile, err = os.Open(filename)
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(ConvertToISO8601(time.Now()), err)
		}

		defer jsonFile.Close()
		byteValue, _ = ioutil.ReadAll(jsonFile)
		if errJson != nil {
			fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
			log.Fatal(errJson)
		}
		errJson := json.Unmarshal(byteValue, &bitQueryResult)
		if errJson != nil {
			fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
			log.Fatal(errJson)
		}
	} else {
		err = ioutil.WriteFile(filename, byteValue, 0644)
		//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", byteValue))
		if err != nil {
			panic(err)
		}
		fmt.Println(ConvertToISO8601(time.Now()), "Succefully written", filename)
	}
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", bitQueryResult.Data.Properties.API))
	var coinData CoinData
	for i := 0; i < len(bitQueryResult.Data.Properties.API); i++ {
		var coinDatum CoinDatum
		coinDatum.BscContract = bitQueryResult.Data.Properties.API[i].Currency.Address
		coinDatum.Properties.Dollar.Price = bitQueryResult.Data.Properties.API[i].Price
		//fmt.Println(ConvertToISO8601(time.Now()),  "contract", coinDatum.BscContract)

		//fmt.Println(ConvertToISO8601(time.Now()),  "coindata from bitquery coindata", fmt.Sprintf("%v", coinDatum))
		coinData.CoinData = append(coinData.CoinData, coinDatum)

	}
	//fmt.Println(ConvertToISO8601(time.Now()),  "coindata from bitquery", fmt.Sprintf("%v", coinData))
	return coinData
}

func ReadBSCPricesFromBitQuery(contract string, contractString string) CoinData {

	//UpdateBSCPricesJsonFromBitQuery(contract, contractString, "BSCPricesTemp.json")
	coinData := ReadBitQueryPriceResultJson(contract, "BSCPricesIn"+contract+".json")

	return coinData
}

func UpdateBSCPricesForContractListJsonFromBitQuery(contract string, contractString string, filename string) {
	yesterday := time.Now()
	yesterday = yesterday.Add(-24 * time.Hour)
	yesterdayString := ConvertToISO8601(yesterday)[:10]
	fmt.Println(yesterdayString)
	query := `{
				  ethereum(network: bsc) {
					dexTrades(
					  options: {desc: ["block.height","tradeIndex"]
					  limitBy: {each: "baseCurrency.address" limit:1}}
					  exchangeName: {is: "Pancake v2"}
					  baseCurrency: {in: ` + contractString + `}
					  quoteCurrency: {is: "` + contract + `"}
					  date: {after: "` + yesterdayString + `"}
					) {

					  tradeIndex
					  block {
						height
						timestamp {
						  iso8601
						}
					  }
					  baseCurrency {
						symbol
						address
					  }
					  quoteCurrency {
						address
					  }
					  quotePrice
				   
					}
				  }
				}
        `
	//fmt.Println(ConvertToISO8601(time.Now()),  query)
	UpdateJsonFromBitQueryFactory(query, filename)
}

func UpdateBSCPricesJsonFromBitQuery(contract string, contractString string, filename string) {
	yesterday := time.Now()
	yesterday = yesterday.Add(-24 * time.Hour)
	yesterdayString := ConvertToISO8601(yesterday)[:10]
	fmt.Println(yesterdayString)
	query := `{
				  ethereum(network: bsc) {
					dexTrades(
					  options: {desc: ["block.height","tradeIndex"]
					  limitBy: {each: "baseCurrency.address" limit:1}}
					  exchangeName: {is: "Pancake v2"}
					  baseCurrency: {not: "` + contract + `"}
					  quoteCurrency: {is: "` + contract + `"}
					  date: {after: "` + yesterdayString + `"}
					  
					) {

					  tradeIndex
					  block {
						height
						timestamp {
						  iso8601
						}
					  }
					  baseCurrency {
						symbol
						address
					  }
					  quoteCurrency {
						address
					  }
					  quotePrice
				   
					}
				  }
				}
        `
	//fmt.Println(ConvertToISO8601(time.Now()),  query)
	UpdateJsonFromBitQueryFactory(query, filename)
}

func UpdateBSCPricesAnyContractJsonFromBitQuery(contractString string, filename string) {
	yesterday := time.Now()
	yesterday = yesterday.Add(-24 * time.Hour)
	yesterdayString := ConvertToISO8601(yesterday)[:10]
	fmt.Println(yesterdayString)
	query := `{
				  ethereum(network: bsc) {
					dexTrades(
					  options: {desc: ["block.height","tradeIndex"]
					  limitBy: {each: "baseCurrency.address" limit:1}}
					  exchangeName: {is: "Pancake v2"}
					  baseCurrency: {in: ` + contractString + `}
					  date: {after: "` + yesterdayString + `"}
					) {
					  tradeIndex
					  block {
						height
						timestamp {
						  iso8601
						}
					  }
					  baseCurrency {
						symbol
						address
					  }
					  quoteCurrency {
						address
					  }
					  quotePrice
				   
					}
				  }
				}
        `
	//fmt.Println(ConvertToISO8601(time.Now()),  query)
	UpdateJsonFromBitQueryFactory(query, filename)
}

func ReadBitQueryQuoteResultFromJson(filename string) BitQueryBSCQuoteResult {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(ConvertToISO8601(time.Now()), err)
	}
	fmt.Println(ConvertToISO8601(time.Now()), "Successfully Opened ", filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	var bitQueryResult BitQueryBSCQuoteResult
	errJson := json.Unmarshal(byteValue, &bitQueryResult)
	if errJson != nil {
		fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
		log.Fatal(errJson)
	}
	return bitQueryResult
}

func RewriteJsonPricesFromBitQueryEvery(duration time.Duration) {
	time.Sleep(20 * time.Second)
	for true {

		var pricePoints BSCPricePoints
		coinDataInDB := ReadCryptosSQLDB("ram")
		contractString := ""
		bscCount := 0
		for i := 0; i < len(coinDataInDB.CoinData); i++ {
			if coinDataInDB.CoinData[i].BscContract != "<nil>" {
				bscCount++

			}

		}
		divider := 20
		sliceSize := math.Ceil(float64(bscCount) / float64(divider))
		dataIndex := 0
		var contractStrings []string
		for k := 0; k < divider; k++ {
			contractString = ""
			countdown := sliceSize
			firstLoop := true
			contractString += "["
			for countdown > 1 && dataIndex < len(coinDataInDB.CoinData) {
				if coinDataInDB.CoinData[dataIndex].BscContract != "<nil>" {
					if !firstLoop {
						contractString += ","

					}
					contractString += "\"" + coinDataInDB.CoinData[dataIndex].BscContract + "\""
					firstLoop = false

					countdown--

				}
				dataIndex++

			}
			contractString += "]"
			fmt.Println(contractString)
			contractStrings = append(contractStrings, contractString)
		}
		for k := 0; k < divider; k++ {
			UpdateBSCPricesForContractListJsonFromBitQuery("0x55d398326f99059ff775485246999027b3197955", contractStrings[k], "AutoBSCPricesUSDT.json")
			if IsPriceJsonFromBitQueryValid("AutoBSCPricesUSDT.json") {
				bitQueryResult := ReadBitQueryQuoteResultFromJson("AutoBSCPricesUSDT.json")
				var coinData1 CoinData
				for i := 0; i < len(bitQueryResult.Data.Properties.API); i++ {
					var coinDatum CoinDatum
					coinDatum.BscContract = bitQueryResult.Data.Properties.API[i].Currency.Address
					coinDatum.Properties.Dollar.Price = bitQueryResult.Data.Properties.API[i].Price
					coinDatum.Properties.Dollar.LastBlock.Height = bitQueryResult.Data.Properties.API[i].Block.Height
					coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = strings.Replace(bitQueryResult.Data.Properties.API[i].Block.TimeStamp.ISO8601, "Z", ".000Z", -1)
					coinDatum.Properties.Dollar.LastBlock.Network = "bsc"
					coinData1.CoinData = append(coinData1.CoinData, coinDatum)
				}
				coinData := ReadCryptosSQLDB("ram")
				coinDataByContract := GetCoinDataByBscContracts(coinData)
				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData1.CoinData[i].Id = coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Id
					if coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name != "" {
						fmt.Println(ConvertToISO8601(time.Now()), "Name", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name, "price_before", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Properties.Dollar.Price, "price_after", bitQueryResult.Data.Properties.API[i].Price, "on_block", coinData1.CoinData[i].Properties.Dollar.LastBlock.Height)
						var PricePoint BSCPricePoint
						PricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
						PricePoint.Block = coinData1.CoinData[i].Properties.Dollar.LastBlock
						PricePoint.CMCId = coinData1.CoinData[i].Id
						pricePoints = append(pricePoints, PricePoint)
					}
				}
				WriteCryptosPriceSQLDB(coinData1, "ram")
				fmt.Println(ConvertToISO8601(time.Now()), "Refreshed prices from USDT", k)
			}
			time.Sleep(3 * time.Second)
		}
		for k := 0; k < divider; k++ {
			UpdateBSCPricesForContractListJsonFromBitQuery("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", contractStrings[k], "AutoBSCPricesBNB.json")
			if IsPriceJsonFromBitQueryValid("AutoBSCPricesBNB.json") {
				bitQueryResult := ReadBitQueryQuoteResultFromJson("AutoBSCPricesBNB.json")
				coinDatum1 := ReadCryptoByBSCContractSQLDB("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", "ram")
				var coinData1 CoinData
				for i := 0; i < len(bitQueryResult.Data.Properties.API); i++ {
					var coinDatum CoinDatum
					coinDatum.BscContract = bitQueryResult.Data.Properties.API[i].Currency.Address
					coinDatum.Properties.Dollar.Price = bitQueryResult.Data.Properties.API[i].Price * coinDatum1.Properties.Dollar.Price
					coinDatum.Properties.Dollar.LastBlock = bitQueryResult.Data.Properties.API[i].Block
					coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = strings.Replace(bitQueryResult.Data.Properties.API[i].Block.TimeStamp.ISO8601, "Z", ".000Z", -1)
					coinDatum.Properties.Dollar.LastBlock.Network = "bsc"
					coinData1.CoinData = append(coinData1.CoinData, coinDatum)

				}
				coinData := ReadCryptosSQLDB("ram")
				coinDataByContract := GetCoinDataByBscContracts(coinData)
				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData1.CoinData[i].Id = coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Id
					if coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name != "" {
						fmt.Println(ConvertToISO8601(time.Now()), "Name", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name, "price_before", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Properties.Dollar.Price, "price_after", coinData1.CoinData[i].Properties.Dollar.Price)
						var PricePoint BSCPricePoint
						PricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
						PricePoint.Block = coinData1.CoinData[i].Properties.Dollar.LastBlock

						PricePoint.CMCId = coinData1.CoinData[i].Id
						pricePoints = append(pricePoints, PricePoint)
					}

				}
				WriteCryptosPriceSQLDB(coinData1, "ram")
				fmt.Println(ConvertToISO8601(time.Now()), "Refreshed prices from BNB", k)
				fmt.Println(coinDatum1.Properties.Dollar.Price)
			}
			time.Sleep(3 * time.Second)
		}

		for k := 0; k < divider; k++ {
			UpdateBSCPricesAnyContractJsonFromBitQuery(contractStrings[k], "AutoBSCPricesANY.json")
			if IsPriceJsonFromBitQueryValid("AutoBSCPricesANY.json") {
				bitQueryResult := ReadBitQueryQuoteResultFromJson("AutoBSCPricesANY.json")
				coinDatum1 := ReadCryptoByBSCContractSQLDB("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", "ram")
				coinData := ReadCryptosSQLDB("ram")
				coinDataByContract := GetCoinDataByBscContracts(coinData)
				for i := 0; i < len(coinData.CoinData); i++ {
					coinDataByContract[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i]
				}
				var coinData1 CoinData
				for i := 0; i < len(bitQueryResult.Data.Properties.API); i++ {
					var coinDatum CoinDatum
					referenceCoinDatum := coinDataByContract[strings.ToLower(bitQueryResult.Data.Properties.API[i].QuoteCurrency.Address)]
					var newTokenDate = strings.Replace(bitQueryResult.Data.Properties.API[i].Block.TimeStamp.ISO8601, "Z", ".000Z", -1)
					var oldTokenDate = referenceCoinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601
					mk := make([]string, 2)
					mk[0] = newTokenDate
					mk[1] = oldTokenDate
					sort.Strings(mk)
					if (referenceCoinDatum.Name != "") && (newTokenDate == mk[0]) {
						coinDatum.BscContract = bitQueryResult.Data.Properties.API[i].Currency.Address
						coinDatum.Properties.Dollar.Price = bitQueryResult.Data.Properties.API[i].Price * referenceCoinDatum.Properties.Dollar.Price
						if bitQueryResult.Data.Properties.API[i].Block.Height < referenceCoinDatum.Properties.Dollar.LastBlock.Height {
							coinDatum.Properties.Dollar.LastBlock.Height = bitQueryResult.Data.Properties.API[i].Block.Height
						} else {
							coinDatum.Properties.Dollar.LastBlock.Height = referenceCoinDatum.Properties.Dollar.LastBlock.Height
						}

						coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = mk[0]
						coinDatum.Properties.Dollar.LastBlock.Network = "bsc"
						coinData1.CoinData = append(coinData1.CoinData, coinDatum)
					}

				}

				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData1.CoinData[i].Id = coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Id
					if coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name != "" {
						fmt.Println(ConvertToISO8601(time.Now()), "Name", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name, "price_before", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Properties.Dollar.Price, "price_after", coinData1.CoinData[i].Properties.Dollar.Price, "on_block", coinData1.CoinData[i].Properties.Dollar.LastBlock.Height)
						var PricePoint BSCPricePoint
						PricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
						PricePoint.Block = coinData1.CoinData[i].Properties.Dollar.LastBlock
						PricePoint.CMCId = coinData1.CoinData[i].Id
						pricePoints = append(pricePoints, PricePoint)
					}

				}
				WriteCryptosPriceSQLDB(coinData1, "ram")
				fmt.Println(ConvertToISO8601(time.Now()), "Refreshed prices from ANY", k)
				fmt.Println(coinDatum1.Properties.Dollar.Price)
			}
			time.Sleep(3 * time.Second)
		}
		fmt.Println(pricePoints)
		WriteBSCPricePointsSQLDB(pricePoints, "ram")
		time.Sleep(duration)
	}
}

func RewriteJsonPricesFromBlockchainEvery(duration time.Duration) {
	time.Sleep(20 * time.Second)
	for true {

		var pricePoints BSCPricePoints
		coinDataInDB := ReadCryptosSQLDB("ram")
		contractString := ""
		bscCount := 0
		for i := 0; i < len(coinDataInDB.CoinData); i++ {
			if coinDataInDB.CoinData[i].BscContract != "<nil>" {
				bscCount++

			}

		}
		divider := 8
		sliceSize := math.Ceil(float64(bscCount) / float64(divider))
		dataIndex := 0
		var contractStrings []string
		for k := 0; k < divider; k++ {
			contractString = ""
			countdown := sliceSize
			firstLoop := true
			for countdown > 1 && dataIndex < len(coinDataInDB.CoinData) {
				if coinDataInDB.CoinData[dataIndex].BscContract != "<nil>" {
					if !firstLoop {
						contractString += "\n"

					}
					contractString += coinDataInDB.CoinData[dataIndex].BscContract
					firstLoop = false

					countdown--

				}
				dataIndex++

			}
			//fmt.Println(contractString)
			contractStrings = append(contractStrings, contractString)
		}
		for k := 0; k < divider; k++ {
			filename := "pyReader/addresses_buffer"
			f, err := os.Create(filename)

			if err != nil {
				log.Fatal(err)
			}
			_, err2 := f.WriteString(contractStrings[k])

			if err2 != nil {
				log.Fatal(err2)
			}
			f.Close()
			LaunchPyReader(filename, "AutoBSCPricesUSDT.json", k)
			if IsPriceJsonFromBitQueryValid("AutoBSCPricesUSDT.json") {
				bitQueryResult := ReadBitQueryQuoteResultFromJson("AutoBSCPricesUSDT.json")
				var coinData1 CoinData
				for i := 0; i < len(bitQueryResult.Data.Properties.API); i++ {
					var coinDatum CoinDatum
					coinDatum.BscContract = bitQueryResult.Data.Properties.API[i].Currency.Address
					coinDatum.Properties.Dollar.Price = bitQueryResult.Data.Properties.API[i].Price
					coinDatum.Properties.Dollar.LastBlock.Height = bitQueryResult.Data.Properties.API[i].Block.Height
					coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = strings.Replace(bitQueryResult.Data.Properties.API[i].Block.TimeStamp.ISO8601, "Z", ".000Z", -1)
					coinDatum.Properties.Dollar.LastBlock.Network = "bsc"
					coinData1.CoinData = append(coinData1.CoinData, coinDatum)
				}
				coinData := ReadCryptosSQLDB("ram")
				coinDataByContract := GetCoinDataByBscContracts(coinData)
				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData1.CoinData[i].Id = coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Id
					if coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name != "" {
						//fmt.Println(ConvertToISO8601(time.Now()), "Name", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name, "price_before", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Properties.Dollar.Price, "price_after", bitQueryResult.Data.Properties.API[i].Price, "on_block", coinData1.CoinData[i].Properties.Dollar.LastBlock.Height)
						var PricePoint BSCPricePoint
						PricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
						PricePoint.Block = coinData1.CoinData[i].Properties.Dollar.LastBlock
						PricePoint.CMCId = coinData1.CoinData[i].Id
						pricePoints = append(pricePoints, PricePoint)
					}
				}
				WriteCryptosPriceSQLDB(coinData1, "ram")
			}
			time.Sleep(3 * time.Second)
		}

		fmt.Println(pricePoints)
		WriteBSCPricePointsSQLDB(pricePoints, "ram")
		time.Sleep(duration)
	}
}

func IsPriceJsonFromBitQueryValid(filename string) bool {

	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(ConvertToISO8601(time.Now()), err)
	}
	fmt.Println(ConvertToISO8601(time.Now()), "Successfully Opened ", filename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var bitQueryResult BitQueryBSCQuoteResult

	errJson := json.Unmarshal(byteValue, &bitQueryResult)
	if errJson != nil {
		fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
		log.Fatal(errJson)
	}

	if len(bitQueryResult.Data.Properties.API) == 0 || bitQueryResult.Data.Properties.API == nil {
		return false
	}
	return true
}

func UpdateJsonFromBitQueryFactory(query string, filename string) {
	client := http.Client{}

	apiKey := ReadApiKey(".env.json").BitQueryApiKey
	//fmt.Println(ConvertToISO8601(time.Now()),  "apikey", apiKey)

	data := map[string]string{
		"query": query,
	}

	jsonValue, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "https://graphql.bitquery.io", bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", err))
		//Handle Error
	}
	//defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(ConvertToISO8601(time.Now()), "Non-OK HTTP status:", res.StatusCode)
		err = ioutil.WriteFile(filename+"trash.json", bodyBytes, 0644)
		if err != nil {
			panic(err)
		}
		// You may read / inspect response body
		return
	} else {
		// write the whole body at once
		err = ioutil.WriteFile(filename, bodyBytes, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println(ConvertToISO8601(time.Now()), "Succefully written", filename)
	}
}
