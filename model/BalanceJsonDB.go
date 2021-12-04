package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ReadBSCScanResultJson(address string, contract string, filename string) BSCBalance {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", filename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var bscScanResult BSCScanBalanceResult

	errJson := json.Unmarshal(byteValue, &bscScanResult)
	if errJson != nil {
		fmt.Println("Failed to unmarshall json file")
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
		UpdateBalanceJson("0xDDd0933873b580313Beb493020F9f72DDA03c9Cb", bscContracts.Contracts[i], "lastbscresult.json")
		bscBalance = ReadBSCScanResultJson("0xDDd0933873b580313Beb493020F9f72DDA03c9Cb", bscContracts.Contracts[i], "lastbscresult.json")
		bscBalances.Balances = append(bscBalances.Balances, bscBalance)
	}
	return bscBalances
}

func UpdateBalanceJson(address string, contract string, filename string) {
	client := http.Client{}

	//apiKey := readApiKey(".env.json").ApiKey
	//fmt.Println("apikey", apiKey)

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
		fmt.Println(fmt.Sprintf("%v", err))
		//Handle Error
	}
	//defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", res.StatusCode)
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
		fmt.Println("Succefully written", filename)
	}
}
