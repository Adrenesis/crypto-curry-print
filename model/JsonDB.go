package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ReadJson(filename string) CoinData {
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
	var coinData CoinData

	errJson := json.Unmarshal(byteValue, &coinData)
	if errJson != nil {
		fmt.Println("Failed to unmarshall json file")
		log.Fatal(errJson)
	}
	return coinData
}

func readApiKey(filename string) ApiKeyData {
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
	var apiKeyData ApiKeyData

	errJson := json.Unmarshal(byteValue, &apiKeyData)
	if errJson != nil {
		fmt.Println("Failed to unmarshall json file")
		log.Fatal(errJson)
	}
	//fmt.Println("apikey", apiKeyData.ApiKey)
	return apiKeyData
}

func updateJson(from string, limit string, filename string) {
	client := http.Client{}

	apiKey := readApiKey(".env.json").ApiKey
	//fmt.Println("apikey", apiKey)

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		//Handle Error
	}
	q := req.URL.Query()
	q.Add("start", from)
	q.Add("limit", limit)
	q.Add("convert", "USD")
	q.Add("sort", "date_added")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		//Handle Error
	}
	defer res.Body.Close()
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

func UpdateJsons(all bool) {
	if all {
		updateJson("1", "5000", "cmcdb0.json")
		updateJson("5001", "5000", "cmcdb1.json")
	} else {
		updateJson("1", "200", "cmcdb200.json")
	}

}
