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
	rows, err := db.Query("select id, name, symbol, price, vol24, date_added, explorer, bscscan, ethscan, xrpscan, bsccontract, ethcontract, xrpcontract from cryptos order by date_added desc;")
	if err != nil {
		log.Fatal(err)
	}
	var coinData CoinData
	for rows.Next() {
		var coinDatum CoinDatum
		var explorer string
		var bscScan interface{}
		var ethScan interface{}
		var xrpScan interface{}
		var bscContract interface{}
		var ethContract interface{}
		var xrpContract interface{}
		if err = rows.Scan(
			&coinDatum.Id,
			&coinDatum.Name,
			&coinDatum.Symbol,
			&coinDatum.Properties.Dollar.Price,
			&coinDatum.Properties.Dollar.Volume24,
			&coinDatum.DateAdded, &explorer,
			&bscScan,
			&ethScan,
			&xrpScan,
			&bscContract,
			&ethContract,
			&xrpContract); err != nil {
			log.Fatal(err)
		}
		coinDatum.Explorers = strings.Split(explorer, ",")
		coinDatum.BscScan = fmt.Sprintf("%v", bscScan)
		coinDatum.EthScan = fmt.Sprintf("%v", ethScan)
		coinDatum.XrpScan = fmt.Sprintf("%v", xrpScan)
		coinDatum.BscContract = fmt.Sprintf("%v", bscContract)
		coinDatum.EthContract = fmt.Sprintf("%v", ethContract)
		coinDatum.XrpContract = fmt.Sprintf("%v", xrpContract)
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
create table if not exists cryptos(id INTEGER, name VARCHAR, symbol VARCHAR, price REAL, vol24 REAL, date_added TEXT, explorer TEXT, bscscan TEXT, ethscan TEXT, xrpscan TEXT, bsccontract TEXT, ethcontract TEXT, xrpcontract TEXT, PRIMARY KEY(id));
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
	stmt, err = db.Prepare(query)
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

func writeBscScan(bscScan string, bscContract string, id int64) {
	db, tx, stmt := Prepare("UPDATE cryptos SET bscscan = ?, bsccontract = ? WHERE id = ?;")
	Exec(tx, stmt, bscScan, bscContract, fmt.Sprintf("%d", id))
	CloseDB(db)
}

func writeEthScan(ethScan string, ethContract string, id int64) {
	db, tx, stmt := Prepare("UPDATE cryptos SET ethscan = ?, ethcontract = ? WHERE id = ?;")
	Exec(tx, stmt, ethScan, ethContract, fmt.Sprintf("%d", id))
	CloseDB(db)
}

func writeXrpScan(xrpScan string, xrpContract string, id int64) {
	db, tx, stmt := Prepare("UPDATE cryptos SET xrpscan = ?, xrpcontract = ? WHERE id = ?;")
	Exec(tx, stmt, xrpScan, xrpContract, fmt.Sprintf("%d", id))
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
		bscScan := ""
		ethScan := ""
		xrpScan := ""
		bscContract := ""
		ethContract := ""
		xrpContract := ""
		for j := 0; j < len(coinDataMap.CoinDataMap[i].URLs.Explorer); j++ {
			if explorer == "" {
				explorer = coinDataMap.CoinDataMap[i].URLs.Explorer[j]

			} else {
				explorer += "," + coinDataMap.CoinDataMap[i].URLs.Explorer[j]
			}
			fmt.Println(coinDataMap.CoinDataMap[i].URLs.Explorer[j])
			fmt.Println(strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "bscscan"))
			if strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "bscscan") {
				bscScan = coinDataMap.CoinDataMap[i].URLs.Explorer[j]

				bscContract = strings.TrimPrefix(bscScan, "https://www.bscscan.com/token/")
				bscContract = strings.TrimPrefix(bscContract, "https://bscscan.com/token/")
				bscContract = strings.TrimPrefix(bscContract, "https://bscscan.com/address/")
				bscContract = strings.TrimPrefix(bscContract, "https://www.bscscan.com/address/")
				if len(bscContract) > 42 {
					bscContract = bscContract[:42]
				}

				if strings.HasPrefix(bscContract, "0x") {
					writeBscScan(bscScan, bscContract, coinDataMap.CoinDataMap[i].Id)
				}

			}
			if strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "etherscan") {
				ethScan = coinDataMap.CoinDataMap[i].URLs.Explorer[j]
				ethContract = strings.TrimPrefix(ethScan, "https://etherscan.io/token/")
				ethContract = strings.TrimPrefix(ethContract, "https://www.etherscan.io/token/")
				writeEthScan(ethScan, ethContract, coinDataMap.CoinDataMap[i].Id)
			}
			if strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "xrpscan") {
				xrpScan = coinDataMap.CoinDataMap[i].URLs.Explorer[j]
				xrpContract = strings.TrimPrefix(xrpScan, "https://xrpscan.com/account/")
				xrpContract = strings.TrimPrefix(xrpContract, "https://www.xrpscan.com/account/")
				writeXrpScan(xrpScan, xrpContract, coinDataMap.CoinDataMap[i].Id)
			}
		}
		//fmt.Println(explorer)

		//args[1] = coinDataMap.CoinDataMap[i].Id
		//args.Explorer = explorer
		//args.Id = coinDataMap.CoinDataMap[i].Id
		writeExplorer(explorer, coinDataMap.CoinDataMap[i].Id)
		//fmt.Println(fmt.Sprintf("%v", res))
	}

}
