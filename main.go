package main

import (
	"encoding/json"
	"fmt"
	"github.com/tyler-sommer/stick"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	//"strconv"
)

type ApiKeyData struct {
	ApiKey string `json:"api_key"`
}

type CoinData struct {
	CoinData []CoinDatum `json:"data"`
}

type CoinDatum struct {
	Name       string   `json:"name"`
	Symbol     string   `json:"symbol"`
	DateAdded  string   `json:"date_added"`
	Properties Property `json:"quote"`
}

type Property struct {
	Dollar MarketValue `json:"USD"`
}
type MarketValue struct {
	Price    float64 `json:"price"`
	Volume24 float64 `json:"volume_24h"`
}

func readJson(filename string) CoinData {
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

func updateJsons() {
	updateJson("1", "5000", "cmcdb0.json")
	updateJson("5001", "5000", "cmcdb1.json")
}

func sortVolumeDecrease(data CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Volume24 > data.CoinData[j].Properties.Dollar.Volume24
	})
}

func sortVolumeIncrease(data CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Volume24 < data.CoinData[j].Properties.Dollar.Volume24
	})
}

func sortPriceDecrease(data CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Price > data.CoinData[j].Properties.Dollar.Price
	})
}

func sortPriceIncrease(data CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Price < data.CoinData[j].Properties.Dollar.Price
	})
}

func main() {
	var coinData CoinData
	var coinData1 CoinData
	coinData = readJson("cmcdb0.json")
	coinData1 = readJson("cmcdb1.json")
	for i := 0; i < len(coinData1.CoinData); i++ {
		coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
	}

	for i := 0; i < len(coinData.CoinData); i++ {
		fmt.Println("Name: ", coinData.CoinData[i].Name)
		fmt.Println("Symbol: ", coinData.CoinData[i].Symbol)
		//var price = coinData.CoinData[i].Properties.Price
		var s = fmt.Sprintf("%.7f", coinData.CoinData[i].Properties.Dollar.Price)
		fmt.Println("Price: ", s)
		s = fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24)
		fmt.Println("Volume 24h: ", s)
		fmt.Println("Date Added: ", coinData.CoinData[i].DateAdded)
	}
	fsRoot, _ := os.Getwd() // Templates are loaded relative to this directory.
	env := stick.New(stick.NewFilesystemLoader(fsRoot))
	env.Filters["number_format"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.10f", v)
	}
	env.Filters["number_format_vol"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.2f", v)
	}
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query()["from"]
		to := r.URL.Query()["to"]
		threshold := r.URL.Query()["price-threshold"]
		submit := r.URL.Query()["submit"]
		incrVol := r.URL.Query()["incr-vol"]
		decrVol := r.URL.Query()["decr-vol"]
		incrPrice := r.URL.Query()["incr-price"]
		decrPrice := r.URL.Query()["decr-price"]
		refresh := r.URL.Query()["refresh"]
		confirm := r.URL.Query()["confirm"]
		fmt.Println("confirm", confirm)
		if (len(refresh) > 0) && (len(confirm) > 0) {
			if confirm[0] == "on" {
				updateJsons()
				coinData = readJson("cmcdb0.json")
				coinData1 = readJson("cmcdb1.json")
				for i := 0; i < len(coinData1.CoinData); i++ {
					coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
				}
			}

		}
		var ifrom, errfrom = strconv.ParseInt(from[0], 10, 64)
		var ito, errto = strconv.ParseInt(to[0], 10, 64)
		var fthresh, errthresh = strconv.ParseFloat(threshold[0], 64)
		if errthresh == nil && errfrom == nil && errto == nil {
			coinData1.CoinData = []CoinDatum{}
			for i := 0; i < len(coinData.CoinData); i++ {
				var i64 = int64(i)
				if ifrom > i64 {
					continue
				}
				if ito < i64 {
					continue
				}
				if coinData.CoinData[i].Properties.Dollar.Price > fthresh {
					continue
				}

				coinData1.CoinData = append(coinData1.CoinData, coinData.CoinData[i])
			}
		} else {
			coinData1.CoinData = coinData.CoinData
		}
		sfrom := "1"
		sto := "10000"
		sthreshold := "0.1"
		if errfrom == nil {
			sfrom = from[0]
		}
		if errto == nil {
			sto = to[0]
		}
		if errthresh == nil {
			sthreshold = threshold[0]
		}
		fmt.Println("Successfully received request", from, to, threshold, len(submit) > 0)
		if len(incrVol) > 0 {
			sortVolumeIncrease(coinData1)
		}
		if len(decrVol) > 0 {
			sortVolumeDecrease(coinData1)
		}
		if len(incrPrice) > 0 {
			sortPriceIncrease(coinData1)
		}
		if len(decrPrice) > 0 {
			sortPriceDecrease(coinData1)
		}

		p := map[string]stick.Value{"coinData": coinData1, "from": sfrom, "to": sto, "threshold": sthreshold}
		var err = env.Execute("index.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
		if err != nil {
			log.Fatal(err)
		}

	})

	http.ListenAndServe(":80", nil)
}
