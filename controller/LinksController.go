package controller

import (
	Model "github.com/Adrenesis/crypto-curry-print/model"
	View "github.com/Adrenesis/crypto-curry-print/view"
	"github.com/tyler-sommer/stick"
	"log"
	"net/http"
	"strconv"
)

func HandleLinks(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()["id"]
	//fmt.Println("confirm", confirm)
	//fmt.Println("refresh", refreshAll)
	var coinDatum Model.CoinDatum
	//fmt.Println((len(refreshAll) > 0) && (len(confirm) > 0))

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
	var iId int64
	if len(id) > 0 {
		iId, _ = strconv.ParseInt(id[0], 10, 64)
	}
	coinDatum = Model.ReadCryptoSQLDB(iId, "ram")
	//fmt.Println(fmt.Sprintf("%v", coinData1))
	env := View.GetEnv()
	//fmt.Println(nil)
	p := map[string]stick.Value{"coinDatum": coinDatum}
	var err = env.Execute("links.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
	if err != nil {
		log.Fatal(err)
	}

}
