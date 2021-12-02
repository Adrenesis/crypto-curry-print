package controller

import (
	Model "../model"
	View "../view"
	"fmt"
	"github.com/tyler-sommer/stick"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func sortVolumeDecrease(data Model.CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Volume24 > data.CoinData[j].Properties.Dollar.Volume24
	})
}

func sortVolumeIncrease(data Model.CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Volume24 < data.CoinData[j].Properties.Dollar.Volume24
	})
}

func sortPriceDecrease(data Model.CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Price > data.CoinData[j].Properties.Dollar.Price
	})
}

func sortPriceIncrease(data Model.CoinData) {
	sort.Slice(data.CoinData[:], func(i, j int) bool {
		return data.CoinData[i].Properties.Dollar.Price < data.CoinData[j].Properties.Dollar.Price
	})
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query()["from"]
	to := r.URL.Query()["to"]
	threshold := r.URL.Query()["price-threshold"]
	submit := r.URL.Query()["submit"]
	incrVol := r.URL.Query()["incr-vol"]
	decrVol := r.URL.Query()["decr-vol"]
	incrPrice := r.URL.Query()["incr-price"]
	decrPrice := r.URL.Query()["decr-price"]
	refresh := r.URL.Query()["refresh"]
	refreshAll := r.URL.Query()["refresh_all"]
	refreshMap := r.URL.Query()["refresh_map"]
	confirm := r.URL.Query()["confirm"]
	fmt.Println("confirm", confirm)
	fmt.Println("refresh", refreshAll)
	var coinData Model.CoinData
	var coinData1 Model.CoinData
	fmt.Println((len(refreshAll) > 0) && (len(confirm) > 0))
	if (len(refresh) > 0) && (len(confirm) > 0) {
		if confirm[0] == "on" {
			Model.UpdateJsons(false)
			coinData = Model.ReadCryptosSQLDB()
			coinData1 = Model.ReadJson("cmcdb200.json")
			for i := 0; i < len(coinData1.CoinData); i++ {
				coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
			}
			Model.WriteCryptosSQLDB(coinData)
		}
	} else if (len(refreshAll) > 0) && (len(confirm) > 0) {

		fmt.Println("getting all cryptocurrencies...")
		if confirm[0] == "on" {
			//Model.UpdateJsons(true)
			coinData = Model.ReadJson("cmcdb0.json")
			coinData1 = Model.ReadJson("cmcdb1.json")
			for i := 0; i < len(coinData1.CoinData); i++ {
				coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
			}
		}
		Model.WriteCryptosSQLDB(coinData)
	} else if (len(refreshMap) > 0) && (len(confirm) > 0) {

		fmt.Println("getting all cryptocurrencies metadata...")
		if confirm[0] == "on" {
			//Model.UpdateMapJsons()
			coinDataMap := Model.ReadMapJsons()
			Model.WriteCryptosMapSQLDB(coinDataMap)
			//coinData = Model.ReadJson("cmcdb0.json")
			//coinData1 = Model.ReadJson("cmcdb1.json")
			//for i := 0; i < len(coinData1.CoinData); i++ {
			//	coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
			//}
		}
		//Model.WriteCryptosSQLDB(coinData)
	} else {

		//coinData = Model.ReadJson("cmcdb0.json")
		//coinData1 = Model.ReadJson("cmcdb1.json")
		//for i := 0; i < len(coinData1.CoinData); i++ {
		//	coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
		//}
		//Model.WriteCryptosSQLDB(coinData)

		coinData = Model.ReadCryptosSQLDB()
	}
	//for i := 0; i < len(coinData.CoinData); i++ {
	//	fmt.Println("Name: ", coinData.CoinData[i].Name)
	//	fmt.Println("Symbol: ", coinData.CoinData[i].Symbol)
	//	//var price = coinData.CoinData[i].Properties.Price
	//	var s = fmt.Sprintf("%.7f", coinData.CoinData[i].Properties.Dollar.Price)
	//	fmt.Println("Price: ", s)
	//	s = fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24)
	//	fmt.Println("Volume 24h: ", s)
	//	fmt.Println("Date Added: ", coinData.CoinData[i].DateAdded)
	//}
	var ifrom int64
	var errfrom error
	if len(from) > 0 {
		ifrom, errfrom = strconv.ParseInt(from[0], 10, 64)
	} else {
		errfrom = strconv.ErrSyntax
	}
	var ito int64
	var errto error
	if len(to) > 0 {
		ito, errto = strconv.ParseInt(to[0], 10, 64)
	} else {
		errto = strconv.ErrSyntax
	}
	var fthresh float64
	var errthresh error
	if len(threshold) > 0 {
		fthresh, errthresh = strconv.ParseFloat(threshold[0], 64)
	} else {
		errthresh = strconv.ErrSyntax
	}
	if errthresh == nil && errfrom == nil && errto == nil {
		coinData1.CoinData = []Model.CoinDatum{}
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
	sfrom := "0"
	sto := "100000"
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

	env := View.GetEnv()
	p := map[string]stick.Value{"coinData": coinData1, "from": sfrom, "to": sto, "threshold": sthreshold}
	var err = env.Execute("index.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
	if err != nil {
		log.Fatal(err)
	}

}
