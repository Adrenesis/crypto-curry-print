package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

func ReadMapJson200() CoinDataMap {
	var coinDataMap CoinDataMap
	coinDataMap1 := ReadMapJson("MapData" + fmt.Sprintf("%d", 200) + ".json")
	//fmt.Println(fmt.Sprintf("%v", coinDataMap1))
	//var elements [][]byte
	for _, element := range coinDataMap1.CoinDataMap {
		//fmt.Println("=>", "Element:", element)
		//for key0, element0 := range element {
		//	fmt.Println("Key:", key0, "=>", "Element:", element0)
		//}

		var coinDatumMap CoinDatumMap

		errJson := json.Unmarshal(element, &coinDatumMap)
		if errJson != nil {
			fmt.Println("Failed to unmarshall json file")
			log.Fatal(errJson)
		}
		coinDataMap.CoinDataMap = append(coinDataMap.CoinDataMap, coinDatumMap)
	}
	return coinDataMap
}

func ReadMapJson(filename string) CoinDataMapUnmarshaler {
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
	var coinDataMap CoinDataMapUnmarshaler

	errJson := json.Unmarshal(byteValue, &coinDataMap)
	if errJson != nil {
		fmt.Println("Failed to unmarshall json file")
		log.Fatal(errJson)
	}
	return coinDataMap
}

func ReadMapJsons() CoinDataMap {
	var coinDataMap CoinDataMap
	for i := 0; i < 16; i++ {
		coinDataMap1 := ReadMapJson("MapData" + fmt.Sprintf("%d", i) + ".json")
		//fmt.Println(fmt.Sprintf("%v", coinDataMap1))
		//var elements [][]byte
		for _, element := range coinDataMap1.CoinDataMap {
			//fmt.Println("=>", "Element:", element)
			//for key0, element0 := range element {
			//	fmt.Println("Key:", key0, "=>", "Element:", element0)
			//}

			var coinDatumMap CoinDatumMap

			errJson := json.Unmarshal(element, &coinDatumMap)
			if errJson != nil {
				fmt.Println("Failed to unmarshall json file")
				log.Fatal(errJson)
			}
			coinDataMap.CoinDataMap = append(coinDataMap.CoinDataMap, coinDatumMap)
		}
		//fmt.Println(fmt.Sprintf("%v", elements))

		//for j := 0; j < len(coinDataMap1.CoinDataMap); j++ {
		//	coinDataMap.CoinDataMap = append(coinDataMap.CoinDataMap, coinDataMap1.CoinDataMap[j])
		//}
	}
	return coinDataMap
}

func ReadApiKey(filename string) ApiKeyData {
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

	apiKey := ReadApiKey(".env.json").CMCApiKey
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
func UpdateMapJsons() {
	coinData := ReadCryptosSQLDB(false)
	client := http.Client{}

	apiKey := ReadApiKey(".env.json").CMCApiKey
	//fmt.Println("apikey", apiKey)

	var ids []int64
	for i := 0; i < len(coinData.CoinData); i++ {
		ids = append(ids, coinData.CoinData[i].Id)
	}
	for j := 0; j < 16; j++ {
		req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/info", nil)
		if err != nil {
			//Handle Error
		}
		var idsString string = ""
		for i := j * 500; i < int(math.Min(float64(len(coinData.CoinData)), float64(499+j*500))); i++ {
			if idsString == "" {
				idsString = fmt.Sprintf("%d", coinData.CoinData[i].Id)
			}
			idsString = idsString + "," + fmt.Sprintf("%d", coinData.CoinData[i].Id)
		}
		q := req.URL.Query()
		q.Add("id", idsString)
		q.Add("aux", "urls,logo,description,tags,platform,date_added,notice,status")
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
			//fmt.Printf("%#v", req.URL.String())
			//fmt.Printf("%#v", req)
			fmt.Printf("%s", bodyBytes)
			// You may read / inspect response body
			return
		} else {
			// write the whole body at once
			filename := "MapData" + fmt.Sprintf("%d", j) + ".json"
			err = ioutil.WriteFile(filename, bodyBytes, 0644)
			if err != nil {
				panic(err)
			}
			fmt.Println("Succefully written", filename)
		}
	}
}

func UpdateMapJson() {
	coinData := ReadCryptosSQLDB(false)
	client := http.Client{}

	apiKey := ReadApiKey(".env.json").CMCApiKey
	//fmt.Println("apikey", apiKey)

	var ids []int64
	for i := 0; i < len(coinData.CoinData); i++ {
		ids = append(ids, coinData.CoinData[i].Id)
	}

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/info", nil)
	if err != nil {
		//Handle Error
	}
	var idsString string = ""
	for i := 0; i < 200; i++ {
		if idsString == "" {
			idsString = fmt.Sprintf("%d", coinData.CoinData[i].Id)
		}
		idsString = idsString + "," + fmt.Sprintf("%d", coinData.CoinData[i].Id)
	}
	q := req.URL.Query()
	q.Add("id", idsString)
	q.Add("aux", "urls,logo,description,tags,platform,date_added,notice,status")
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
		//fmt.Printf("%#v", req.URL.String())
		//fmt.Printf("%#v", req)
		fmt.Printf("%s", bodyBytes)
		// You may read / inspect response body
		return
	} else {
		// write the whole body at once
		filename := "MapData" + fmt.Sprintf("%d", 200) + ".json"
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
