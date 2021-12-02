package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	//"time"
)

func ReadCryptosSQLDB() CoinData {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}
	CreateTable()
	rows, err := db.Query("select id, name, symbol, price, vol24, date_added, explorer from cryptos order by date_added desc;")
	if err != nil {
		log.Fatal(err)
	}
	var coinData CoinData
	for rows.Next() {
		var coinDatum CoinDatum
		if err = rows.Scan(&coinDatum.Id, &coinDatum.Name, &coinDatum.Symbol, &coinDatum.Properties.Dollar.Price, &coinDatum.Properties.Dollar.Volume24, &coinDatum.DateAdded, &coinDatum.Explorers); err != nil {
			log.Fatal(err)
		}
		coinData.CoinData = append(coinData.CoinData, coinDatum)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(fmt.Sprintf("%v", coinData))

	CloseDB(db)
	return coinData
}

func CloseDB(db *sql.DB) {

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

}

func CreateTable() {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists cryptos(id INTEGER, name VARCHAR, symbol VARCHAR, price REAL, vol24 REAL, date_added TEXT, explorer TEXT, PRIMARY KEY(id));
	`); err != nil {
		log.Fatal(err)
	}
}

func WriteCryptosSQLDB(coinData CoinData) {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	CreateTable()

	for i := 0; i < len(coinData.CoinData); i++ {
		//fmt.Println("INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
		//	strings.Replace(coinData.CoinData[i].Name, "'", "''",-1) +"', '" +
		//	strings.Replace(coinData.CoinData[i].Symbol, "'", "''",-1) +"', '" +
		//	fmt.Sprintf("%.7+f", coinData.CoinData[i].Properties.Dollar.Price) +"', '" +
		//	fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) +"', '" +
		//	coinData.CoinData[i].DateAdded +"');")
		if _, err = db.Exec("INSERT INTO cryptos (id, name, symbol, price, vol24, date_added, explorer) VALUES ('" +
			fmt.Sprintf("%d", coinData.CoinData[i].Id) + "', '" +
			strings.Replace(coinData.CoinData[i].Name, "'", "''", -1) + "', '" +
			strings.Replace(coinData.CoinData[i].Symbol, "'", "''", -1) + "', '" +
			fmt.Sprintf("%.7f", coinData.CoinData[i].Properties.Dollar.Price) + "', '" +
			fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) + "', '" +
			coinData.CoinData[i].DateAdded + "', '" +
			"" + "');"); err != nil {
			s := fmt.Sprintf("%v", err)
			if !strings.HasSuffix(s, "(1555)") {
				log.Fatal(err)
			}
		}
	}
	CloseDB(db)

}

func Prepare(query string) (db *sql.DB, tx *sql.Tx, stmt *sql.Stmt) {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err = db.Prepare("UPDATE cryptos SET explorer = ? WHERE id = ?;")
	if err != nil {
		log.Fatal(err)

	}

	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)

	}

	return db, tx, stmt
}

func Exec(tx *sql.Tx, stmt *sql.Stmt, args ...interface{}) {
	defer tx.Rollback()
	_, err := stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)

	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func writeExplorer(explorer string, id int64) {
	db, tx, updateExplorers := Prepare("UPDATE cryptos SET explorer = ? WHERE id = ?;")
	Exec(tx, updateExplorers, explorer, fmt.Sprintf("%d", id))
	CloseDB(db)
}

func WriteCryptosMapSQLDB(coinDataMap CoinDataMap) {

	CreateTable()

	for i := 0; i < len(coinDataMap.CoinDataMap); i++ {
		//fmt.Println("INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
		//	strings.Replace(coinData.CoinData[i].Name, "'", "''",-1) +"', '" +
		//	strings.Replace(coinData.CoinData[i].Symbol, "'", "''",-1) +"', '" +
		//	fmt.Sprintf("%.7+f", coinData.CoinData[i].Properties.Dollar.Price) +"', '" +
		//	fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) +"', '" +
		//	coinData.CoinData[i].DateAdded +"');")
		explorer := ""
		for j := 0; j < len(coinDataMap.CoinDataMap[i].URLs.Explorer); j++ {
			if explorer == "" {
				explorer = coinDataMap.CoinDataMap[i].URLs.Explorer[j]
			} else {
				explorer += "," + coinDataMap.CoinDataMap[i].URLs.Explorer[j]
			}
		}
		fmt.Println(explorer)

		//args[1] = coinDataMap.CoinDataMap[i].Id
		//args.Explorer = explorer
		//args.Id = coinDataMap.CoinDataMap[i].Id
		writeExplorer(explorer, coinDataMap.CoinDataMap[i].Id)
		//fmt.Println(fmt.Sprintf("%v", res))
	}

}
