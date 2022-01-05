package controller

import (
	Model "github.com/Adrenesis/crypto-curry-print/model"
	View "github.com/Adrenesis/crypto-curry-print/view"
	"github.com/tyler-sommer/stick"
	"log"
	"net/http"
	"strconv"
)

func HandleCoinUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()["id"]
	isReviewed := r.URL.Query()["isReviewed"]
	isInteresting := r.URL.Query()["isInteresting"]
	isRejected := r.URL.Query()["isRejected"]

	var coinDatum Model.CoinDatum

	var err error
	var iId int64
	if len(id) > 0 {
		iId, err = strconv.ParseInt(id[0], 10, 64)
	}
	coinDatum = Model.ReadCryptoSQLDB(iId, "ram")
	status := http.StatusOK
	isSomethingUpdated := false
	if len(isReviewed) > 0 {
		coinDatum.IsReviewed, err = strconv.ParseBool(isReviewed[0])
		if err != nil {
			status = http.StatusInternalServerError
		} else {
			isSomethingUpdated = true
		}
	}
	if len(isInteresting) > 0 {
		coinDatum.IsInteresting, err = strconv.ParseBool(isInteresting[0])
		if err != nil {
			status = http.StatusInternalServerError
		} else {
			isSomethingUpdated = true
		}

	}
	if len(isRejected) > 0 {
		coinDatum.IsRejected, err = strconv.ParseBool(isRejected[0])
		if err != nil {
			status = http.StatusInternalServerError
		} else {
			isSomethingUpdated = true
		}
	}
	if !isSomethingUpdated {
		status = http.StatusInternalServerError
	}
	var coinData Model.CoinData
	var coinDatum1 Model.CoinDatum
	coinData.CoinData = append(coinData.CoinData, coinDatum)
	if status == http.StatusOK {
		Model.WriteCryptosFullSQLDB(coinData, "ram")
		coinDatum1 = Model.ReadCryptoSQLDB(iId, "ram")
	}
	env := View.GetEnv()
	p := map[string]stick.Value{"status": status, "coinDatum": coinDatum1}

	err = env.Execute("json.confirm_update.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
	if err != nil {
		log.Fatal(err)
	}

}
