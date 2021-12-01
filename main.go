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
	fmt.Println("Successfully Opened testcmc.json")
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
	//// Open our jsonFile
	//jsonFile, err := os.Open("testcmc.json")
	//// if we os.Open returns an error then handle it
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Successfully Opened testcmc.json")
	//// defer the closing of our jsonFile so that we can parse it later on
	//defer jsonFile.Close()
	//
	//// read our opened xmlFile as a byte array.
	//byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var coinData CoinData
	var coinData1 CoinData
	coinData = readJson("testcmc.json")
	coinData1 = readJson("testcmc1.json")
	for i := 0; i < len(coinData1.CoinData); i++ {
		coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
	}
	//coinData = append(coinData, []CoinData coinData1)
	//fmt.Println("Name: ")
	//errJson := json.Unmarshal(byteValue, &coinData)
	//if errJson != nil {
	//	fmt.Println("Failed to unmarshall json file")
	//	log.Fatal(errJson)
	//	return
	//}

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
		return fmt.Sprintf("%.7f", v)
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
		//sortVolumeDecrease(coinData1)
		p := map[string]stick.Value{"coinData": coinData1, "from": sfrom, "to": sto, "threshold": sthreshold}
		var err = env.Execute("index.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
		if err != nil {
			log.Fatal(err)
		}

	})

	http.ListenAndServe(":80", nil)
}
