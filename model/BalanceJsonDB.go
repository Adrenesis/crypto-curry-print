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
	"strconv"
	"strings"
	"time"
)

func ReadBSCScanResultJson(address string, contract string, filename string) BSCBalance {
	// Open our jsonFile
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
	var bscScanResult BSCScanBalanceResult

	errJson := json.Unmarshal(byteValue, &bscScanResult)
	if errJson != nil {
		fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
		log.Fatal(errJson)
	}
	var bscBalance BSCBalance
	balance, _ := strconv.ParseFloat(bscScanResult.Amount, 64)
	bscBalance.Amount = balance / 1000000000000000000.0
	bscBalance.Address = address
	bscBalance.Contract = contract
	return bscBalance
}

func ReadBSCBalancesFromBSCScan(bscContracts BSCContracts) BSCBalances {
	var bscBalances BSCBalances
	for i := 0; i < 1; i++ {
		var bscBalance BSCBalance
		UpdateBSCScanBalanceJson("0xDDd0933873b580313Beb493020F9f72DDA03c9Cb", bscContracts.Contracts[i], "lastbscresult.json")
		bscBalance = ReadBSCScanResultJson("0xDDd0933873b580313Beb493020F9f72DDA03c9Cb", bscContracts.Contracts[i], "lastbscresult.json")
		bscBalances.Balances = append(bscBalances.Balances, bscBalance)
	}
	return bscBalances
}

func UpdateBSCScanBalanceJson(address string, contract string, filename string) {
	client := http.Client{}

	//apiKey := readApiKey(".env.json").ApiKey
	//fmt.Println(ConvertToISO8601(time.Now()),  "apikey", apiKey)

	req, err := http.NewRequest("GET", "https://api.bscscan.com/api", nil)
	if err != nil {
		//Handle Error
	}
	q := req.URL.Query()
	q.Add("module", "account")
	q.Add("action", "tokenbalance")
	q.Add("contractaddress", contract)
	q.Add("address", address)
	q.Add("tag", "latest")
	q.Add("apikey", ReadApiKey(".env.json").BSCApiKey)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", err))
		//Handle Error
	}
	//defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(ConvertToISO8601(time.Now()), "Non-OK HTTP status:", res.StatusCode)
		fmt.Printf("%#v", req.URL.String())
		fmt.Printf("%#v", req)
		fmt.Printf("%#v", q)
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

func ReadBitQueryResultJson(address string, filename string) BSCBalances {
	// Open our jsonFile
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
	var bitQueryResult BitQueryBSCBalanceResult

	errJson := json.Unmarshal(byteValue, &bitQueryResult)
	if errJson != nil {
		fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
		log.Fatal(errJson)
	}
	var bscBalances BSCBalances
	//fmt.Println(ConvertToISO8601(time.Now()),  "bitquery result", fmt.Sprintf("%v", bitQueryResult))
	for i := 0; i < len(bitQueryResult.Data.Properties.API[0].Balances); i++ {
		fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", bitQueryResult.Data.Properties.API[0].Balances[i].Token.Address))
		var bscBalance BSCBalance
		bscBalance.Amount = bitQueryResult.Data.Properties.API[0].Balances[i].Amount
		bscBalance.Address = address
		bscBalance.Contract = bitQueryResult.Data.Properties.API[0].Balances[i].Token.Address
		if bscBalance.Amount > 0 {
			bscBalances.Balances = append(bscBalances.Balances, bscBalance)
		}
	}

	return bscBalances
}

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

func ReadBSCBalancesFromBitQuery(address string) BSCBalances {

	UpdateBSCBalanceJsonFromBitQuery(address, address+"bscbalance.json")
	bscBalances := ReadBitQueryResultJson(address, address+"bscbalance.json")

	return bscBalances
}

func UpdateBSCBalanceJsonFromBitQuery(address string, filename string) {
	query := `
			{
			  ethereum(network: bsc) {
				address(address: {is: "` + address + `"}) {
				  balances {
					value
					currency {
					  symbol
					  address
					}
				  }
				}
			  }
			}
        `
	UpdateJsonFromBitQueryFactory(query, filename)
}

func UpdateBSCPricesForContractListJsonFromBitQuery(contract string, contractString string, filename string) {
	query := `{
				  ethereum(network: bsc) {
					dexTrades(
					  options: {desc: ["block.height","tradeIndex"]
					  limitBy: {each: "baseCurrency.address" limit:1}}
					  exchangeName: {is: "Pancake v2"}
					  baseCurrency: {in: ` + contractString + `}
					  quoteCurrency: {is: "` + contract + `"}
					  
					) {

					  tradeIndex
					  block {
						height
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
	query := `{
				  ethereum(network: bsc) {
					dexTrades(
					  options: {desc: ["block.height","tradeIndex"]
					  limitBy: {each: "baseCurrency.address" limit:1}}
					  exchangeName: {is: "Pancake v2"}
					  baseCurrency: {not: "` + contract + `"}
					  quoteCurrency: {is: "` + contract + `"}
					  
					) {

					  tradeIndex
					  block {
						height
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
	query := `{
				  ethereum(network: bsc) {
					dexTrades(
					  options: {desc: ["block.height","tradeIndex"]
					  limitBy: {each: "baseCurrency.address" limit:1}}
					  exchangeName: {is: "Pancake v2"}
					  baseCurrency: {not: "0x55d398326f99059ff775485246999027b3197955"}
					) {
					  tradeIndex
					  block {
						height
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
	for true {
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
					coinData1.CoinData = append(coinData1.CoinData, coinDatum)
				}
				coinData := ReadCryptosSQLDB("ram")
				coinDataIndexByContract := make(map[string]int64)
				coinDataByContract := make(map[string]CoinDatum)
				for i := 0; i < len(coinData.CoinData); i++ {
					coinDataIndexByContract[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i].Id
					coinDataByContract[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i]
				}
				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData1.CoinData[i].Id = coinDataIndexByContract[strings.ToLower(coinData1.CoinData[i].BscContract)]
					fmt.Println(ConvertToISO8601(time.Now()), "Name", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name, "price_before", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Properties.Dollar.Price, "price_after", bitQueryResult.Data.Properties.API[i].Price)
				}
				WriteCryptosPriceSQLDB(coinData1, "ram")
				fmt.Println(ConvertToISO8601(time.Now()), "Refreshed prices from USDT", k)
			}
		}
		fmt.Println(contractStrings[0])
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
					coinData1.CoinData = append(coinData1.CoinData, coinDatum)

				}
				coinData := ReadCryptosSQLDB("ram")
				coinDataIndex := make(map[string]int64)
				coinDataByContract := make(map[string]CoinDatum)
				for i := 0; i < len(coinData.CoinData); i++ {
					coinDataIndex[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i].Id
					coinDataByContract[strings.ToLower(coinData.CoinData[i].BscContract)] = coinData.CoinData[i]
				}
				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData1.CoinData[i].Id = coinDataIndex[strings.ToLower(coinData1.CoinData[i].BscContract)]
					if coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name != "" {
						fmt.Println(ConvertToISO8601(time.Now()), "Name", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Name, "price_before", coinDataByContract[strings.ToLower(coinData1.CoinData[i].BscContract)].Properties.Dollar.Price, "price_after", coinData1.CoinData[i].Properties.Dollar.Price)
					}

				}
				WriteCryptosPriceSQLDB(coinData1, "ram")
				fmt.Println(ConvertToISO8601(time.Now()), "Refreshed prices from BNB", k)
				fmt.Println(coinDatum1.Properties.Dollar.Price)
			}
		}
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

func ReadUnmarshalResultJson(address string, filename string) BSCBalances {
	// Open our jsonFile
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
	var unmarshalResult UnmarshalBSCBalanceResult

	errJson := json.Unmarshal(byteValue, &unmarshalResult)
	if errJson != nil {
		fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
	}
	var bscBalances BSCBalances
	//fmt.Println(ConvertToISO8601(time.Now()),  "unmarshal result", fmt.Sprintf("%v", unmarshalResult))
	for i := 0; i < len(unmarshalResult); i++ {
		fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", unmarshalResult[i]))
		var bscBalance BSCBalance
		amount, _ := strconv.ParseFloat(unmarshalResult[i].Amount, 64)
		bscBalance.Amount = amount / math.Pow(10, unmarshalResult[i].Decimals)
		bscBalance.Address = address
		bscBalance.Contract = unmarshalResult[i].Contract
		if bscBalance.Amount > 0 {
			bscBalances.Balances = append(bscBalances.Balances, bscBalance)
		}
	}
	return bscBalances
}

func UpdateJsonFromUnmarshal(address string, filename string) {
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://stg-api.unmarshal.io/v1/bsc/address/"+address+"/assets", nil)

	//req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", err))
		//Handle Error
	}
	//defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(ConvertToISO8601(time.Now()), "Non-OK HTTP status:", res.StatusCode)
		// You may read / inspect response body
		return
	} else {
		// we initialize our Users array
		var unmarshalResult UnmarshalBSCBalanceResult

		errJson := json.Unmarshal(bodyBytes, &unmarshalResult)
		if errJson != nil {
			fmt.Println(ConvertToISO8601(time.Now()), "Failed to unmarshall json file")
			//log.Fatal(errJson)

			return
		}
		// write the whole body at once
		err = ioutil.WriteFile(filename, bodyBytes, 0644)
		//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", bodyBytes))
		if err != nil {
			panic(err)
		}
		fmt.Println(ConvertToISO8601(time.Now()), "Succefully written", filename)
	}
}

func ReadBSCBalancesFromUnmarshal(address string) BSCBalances {
	fmt.Println(ConvertToISO8601(time.Now()), "update unmarshal json...")
	UpdateJsonFromUnmarshal(address, "BSCBalancesFromUnmarshal"+address+".json")
	fmt.Println(ConvertToISO8601(time.Now()), "read unmarshal json...")
	bscBalances := ReadUnmarshalResultJson(address, "BSCBalancesFromUnmarshal"+address+".json")

	return bscBalances
}
