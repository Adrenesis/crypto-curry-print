package controller

import (
	"fmt"
	Model "github.com/Adrenesis/crypto-curry-print/model"
	View "github.com/Adrenesis/crypto-curry-print/view"
	"github.com/tyler-sommer/stick"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var DBinited = false
var cData Model.CoinData

func InitDB() {
	DBinited = true
	Model.CreateBSCBalancesTable("hdd")
	Model.CreateBSCContractsTable("hdd")
	Model.CreateCryptoTable("hdd")
	Model.CreateBSCaddressesTable("hdd")
	Model.CreateBSCPricePointTable("hdd")
	Model.CreateBSCBalancesTable("ram")
	Model.CreateBSCContractsTable("ram")
	Model.CreateBSCaddressesTable("ram")
	Model.CreateCryptoTable("ram")
	Model.CreateBSCPricePointTable("ram")
	cData = Model.ReadCryptosSQLDB("hdd")
	bscBalances := Model.ReadBSCBalancesSQLDB("hdd")
	bscContracts := Model.ReadBSCContractsQLDB("hdd")
	bscAddresses := Model.ReadBSCaddressesSQLDB("hdd")
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", cData))
	//Model.CreateCryptoTable()
	//Model.WriteCryptosSQLDB(cData)
	Model.WriteCryptosFullSQLDB(cData, "ram")
	//var cd Model.CoinData
	//cd = Model.ReadCryptosSQLDB(false)
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", cd))
	//fmt.Println(ConvertToISO8601(time.Now()),  "This is from RAM")
	//cData = Model.ReadCryptosSQLDB()
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", cData))
	Model.WriteBSCContractsSQLDB(bscContracts, "ram")
	Model.WriteBSCBalancesSQLDB(bscBalances, "ram")
	Model.WriteBSCaddressesSQLDB(bscAddresses, "ram")
}

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
	//fmt.Println(ConvertToISO8601(time.Now()),  "confirm", confirm)
	//fmt.Println(ConvertToISO8601(time.Now()),  "refresh", refreshAll)
	var coinData Model.CoinData
	var coinData1 Model.CoinData
	//fmt.Println(ConvertToISO8601(time.Now()),  (len(refreshAll) > 0) && (len(confirm) > 0))
	if (len(refresh) > 0) && (len(confirm) > 0) {
		if confirm[0] == "on" {
			Model.UpdateJsons(false)
			coinData = Model.ReadCryptosSQLDB("ram")
			coinData1 = Model.ReadJson("cmcdb200.json")
			var pricePoints Model.BSCPricePoints
			for i := 0; i < len(coinData1.CoinData); i++ {
				coinData1.CoinData[i].Properties.Dollar.LastBlock.Network = "cmc"
				coinData1.CoinData[i].Properties.Dollar.LastBlock.Height = -1
				coinData1.CoinData[i].Properties.Dollar.LastBlock.TimeStamp.ISO8601 = coinData1.Status.TimeStamp
				coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
				var pricePoint Model.BSCPricePoint
				pricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
				pricePoint.Block.Height = -1
				pricePoint.Block.TimeStamp.ISO8601 = coinData1.Status.TimeStamp
				pricePoint.CMCId = coinData1.CoinData[i].Id
				pricePoints = append(pricePoints, pricePoint)
			}
			Model.WriteCryptosSQLDB(coinData, "ram")
			Model.UpdateMapJson()
			var coinDataMap Model.CoinDataMap
			coinDataMap = Model.ReadMapJson200()
			Model.WriteCryptosMapSQLDB(coinDataMap, "ram")
			Model.WriteBSCPricePointsSQLDB(pricePoints, "ram")
		}
		coinData = Model.ReadCryptosSQLDB("ram")
	} else if (len(refreshAll) > 0) && (len(confirm) > 0) {

		fmt.Println(Model.ConvertToISO8601(time.Now()), "getting all cryptocurrencies...")
		if confirm[0] == "on" {
			Model.UpdateJsons(true)

			var pricePoints Model.BSCPricePoints
			coinData1 = Model.ReadJson("cmcdb0.json")
			for i := 0; i < len(coinData1.CoinData); i++ {
				coinData1.CoinData[i].Properties.Dollar.LastBlock.Network = "cmc"
				coinData1.CoinData[i].Properties.Dollar.LastBlock.Height = -1
				coinData1.CoinData[i].Properties.Dollar.LastBlock.TimeStamp.ISO8601 = coinData1.Status.TimeStamp
				coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
				var pricePoint Model.BSCPricePoint
				pricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
				pricePoint.Block.Height = -1
				pricePoint.Block.TimeStamp.ISO8601 = coinData1.Status.TimeStamp
				pricePoint.CMCId = coinData1.CoinData[i].Id
				pricePoints = append(pricePoints, pricePoint)
			}

			coinData1 = Model.ReadJson("cmcdb1.json")
			for i := 0; i < len(coinData1.CoinData); i++ {
				coinData1.CoinData[i].Properties.Dollar.LastBlock.Network = "cmc"
				coinData1.CoinData[i].Properties.Dollar.LastBlock.Height = -1
				coinData1.CoinData[i].Properties.Dollar.LastBlock.TimeStamp.ISO8601 = coinData1.Status.TimeStamp
				coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
				//if coinData1.CoinData[i].BscContract != "" {
				var pricePoint Model.BSCPricePoint
				pricePoint.Price = coinData1.CoinData[i].Properties.Dollar.Price
				pricePoint.Block.Height = -1
				pricePoint.Block.TimeStamp.ISO8601 = coinData1.Status.TimeStamp
				pricePoint.CMCId = coinData1.CoinData[i].Id
				pricePoints = append(pricePoints, pricePoint)
				//}
			}
			fmt.Println(pricePoints)
			Model.WriteBSCPricePointsSQLDB(pricePoints, "ram")
		}
		Model.WriteCryptosSQLDB(coinData, "ram")
		coinData = Model.ReadCryptosSQLDB("ram")
	} else if (len(refreshMap) > 0) && (len(confirm) > 0) {

		fmt.Println(Model.ConvertToISO8601(time.Now()), "getting all cryptocurrencies metadata...")
		if confirm[0] == "on" {
			Model.UpdateMapJsons()
			coinDataMap := Model.ReadMapJsons()
			Model.WriteCryptosMapSQLDB(coinDataMap, "ram")
			//coinData = Model.ReadJson("cmcdb0.json")
			//coinData1 = Model.ReadJson("cmcdb1.json")
			//for i := 0; i < len(coinData1.CoinData); i++ {
			//	coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
			//}
			coinData = Model.ReadCryptosSQLDB("ram")
		}
		//Model.WriteCryptosSQLDB(coinData)
	} else {

		//coinData = Model.ReadJson("cmcdb0.json")
		//coinData1 = Model.ReadJson("cmcdb1.json")
		//for i := 0; i < len(coinData1.CoinData); i++ {
		//	coinData.CoinData = append(coinData.CoinData, coinData1.CoinData[i])
		//}
		//Model.WriteCryptosSQLDB(coinData)

		coinData = Model.ReadCryptosSQLDB("ram")
	}
	//for i := 0; i < len(coinData.CoinData); i++ {
	//	fmt.Println(ConvertToISO8601(time.Now()),  "Name: ", coinData.CoinData[i].Name)
	//	fmt.Println(ConvertToISO8601(time.Now()),  "Symbol: ", coinData.CoinData[i].Symbol)
	//	//var price = coinData.CoinData[i].Properties.Price
	//	var s = fmt.Sprintf("%.7f", coinData.CoinData[i].Properties.Dollar.Price)
	//	fmt.Println(ConvertToISO8601(time.Now()),  "Price: ", s)
	//	s = fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24)
	//	fmt.Println(ConvertToISO8601(time.Now()),  "Volume 24h: ", s)
	//	fmt.Println(ConvertToISO8601(time.Now()),  "Date Added: ", coinData.CoinData[i].DateAdded)
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
	fmt.Println(Model.ConvertToISO8601(time.Now()), "Successfully received request", from, to, threshold, len(submit) > 0)
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
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", coinData1))
	env := View.GetEnv()
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("#%v", Model.DBReady))
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("#%v", Model.DBAlt))
	//fmt.Println(ConvertToISO8601(time.Now()),  "db inited", DBinited)
	p := map[string]stick.Value{"coinData": coinData1, "from": sfrom, "to": sto, "threshold": sthreshold}
	var err = env.Execute("index.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
	if err != nil {
		log.Fatal(err)
	}

}
