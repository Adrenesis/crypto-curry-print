package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
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
		//fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", bitQueryResult.Data.Properties.API[0].Balances[i].Token.Address))
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
		//fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("%v", unmarshalResult[i]))
		var bscBalance BSCBalance
		amount, _ := strconv.ParseFloat(unmarshalResult[i].Amount, 64)
		bscBalance.Amount = amount / math.Pow(10, unmarshalResult[i].Decimals)
		bscBalance.Address = address
		if unmarshalResult[i].Contract == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
			bscBalance.Contract = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
		} else {
			bscBalance.Contract = unmarshalResult[i].Contract
		}

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
