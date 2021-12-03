package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	//"time"
)

func ReadCryptoSQLDB(id int64) CoinDatum {
	fmt.Println("reading database...")
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, name, slug, symbol, logo, price, vol24, date_added, explorer, bscscan, ethscan, xrpscan, bsccontract, ethcontract, xrpcontract, twitter, website, facebook, chat, message_board, technical, source_code, announcement from cryptos WHERE id = ?  order by date_added desc;", fmt.Sprintf("%d", id))
	if err != nil {
		log.Fatal(err)
	}

	var coinDatum CoinDatum
	for rows.Next() {
		var slug interface{}
		var logo interface{}
		var explorer string
		var twitter interface{}
		var website interface{}
		var facebook interface{}
		var chat interface{}
		var messageBoard interface{}
		var technical interface{}
		var sourceCode interface{}
		var announcement interface{}
		var bscScan interface{}
		var ethScan interface{}
		var xrpScan interface{}
		var bscContract interface{}
		var ethContract interface{}
		var xrpContract interface{}
		if err = rows.Scan(
			&coinDatum.Id,
			&coinDatum.Name,
			&slug,
			&coinDatum.Symbol,
			&logo,
			&coinDatum.Properties.Dollar.Price,
			&coinDatum.Properties.Dollar.Volume24,
			&coinDatum.DateAdded,
			&explorer,
			&bscScan,
			&ethScan,
			&xrpScan,
			&bscContract,
			&ethContract,
			&xrpContract,
			&twitter,
			&website,
			&facebook,
			&chat,
			&messageBoard,
			&technical,
			&sourceCode,
			&announcement); err != nil {
			log.Fatal(err)
		}
		coinDatum.Explorers = strings.Split(explorer, ",")
		coinDatum.Twitters = strings.Split(fmt.Sprintf("%v", twitter), ",")
		coinDatum.Facebooks = strings.Split(fmt.Sprintf("%v", facebook), ",")
		coinDatum.Websites = strings.Split(fmt.Sprintf("%v", website), ",")
		coinDatum.MessageBoards = strings.Split(fmt.Sprintf("%v", messageBoard), ",")
		coinDatum.Chats = strings.Split(fmt.Sprintf("%v", chat), ",")
		coinDatum.Technicals = strings.Split(fmt.Sprintf("%v", technical), ",")
		coinDatum.SourceCodes = strings.Split(fmt.Sprintf("%v", sourceCode), ",")
		coinDatum.Announcements = strings.Split(fmt.Sprintf("%v", announcement), ",")
		coinDatum.Slug = fmt.Sprintf("%v", slug)
		if coinDatum.Slug == "<nil>" {
			coinDatum.Slug = strings.ToLower(coinDatum.Name)
			coinDatum.Slug = strings.Replace(coinDatum.Slug, " ", "-", -1)
		}
		coinDatum.Logo = fmt.Sprintf("%v", logo)
		//fmt.Println("logo", coinDatum.Logo)
		coinDatum.BscScan = fmt.Sprintf("%v", bscScan)
		coinDatum.EthScan = fmt.Sprintf("%v", ethScan)
		coinDatum.XrpScan = fmt.Sprintf("%v", xrpScan)
		coinDatum.BscContract = fmt.Sprintf("%v", bscContract)
		coinDatum.EthContract = fmt.Sprintf("%v", ethContract)
		coinDatum.XrpContract = fmt.Sprintf("%v", xrpContract)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(fmt.Sprintf("%v", coinData))

	CloseDB(db)
	return coinDatum
}

func ReadCryptosSQLDB() CoinData {
	fmt.Println("reading database...")
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}
	CreateTable()
	rows, err := db.Query("select id, name, slug, symbol, logo, price, vol24, date_added, explorer, bscscan, ethscan, xrpscan, bsccontract, ethcontract, xrpcontract, twitter, website, facebook, chat, message_board, technical, source_code, announcement from cryptos order by date_added desc;")
	if err != nil {
		log.Fatal(err)
	}
	var coinData CoinData
	for rows.Next() {
		var coinDatum CoinDatum
		var slug interface{}
		var logo interface{}
		var explorer string
		var twitter interface{}
		var website interface{}
		var facebook interface{}
		var chat interface{}
		var messageBoard interface{}
		var technical interface{}
		var sourceCode interface{}
		var announcement interface{}
		var bscScan interface{}
		var ethScan interface{}
		var xrpScan interface{}
		var bscContract interface{}
		var ethContract interface{}
		var xrpContract interface{}
		if err = rows.Scan(
			&coinDatum.Id,
			&coinDatum.Name,
			&slug,
			&coinDatum.Symbol,
			&logo,
			&coinDatum.Properties.Dollar.Price,
			&coinDatum.Properties.Dollar.Volume24,
			&coinDatum.DateAdded,
			&explorer,
			&bscScan,
			&ethScan,
			&xrpScan,
			&bscContract,
			&ethContract,
			&xrpContract,
			&twitter,
			&website,
			&facebook,
			&chat,
			&messageBoard,
			&technical,
			&sourceCode,
			&announcement); err != nil {
			log.Fatal(err)
		}
		coinDatum.Explorers = strings.Split(explorer, ",")
		coinDatum.Twitters = strings.Split(fmt.Sprintf("%v", twitter), ",")
		coinDatum.Facebooks = strings.Split(fmt.Sprintf("%v", facebook), ",")
		coinDatum.Websites = strings.Split(fmt.Sprintf("%v", website), ",")
		coinDatum.MessageBoards = strings.Split(fmt.Sprintf("%v", messageBoard), ",")
		coinDatum.Chats = strings.Split(fmt.Sprintf("%v", chat), ",")
		coinDatum.Technicals = strings.Split(fmt.Sprintf("%v", technical), ",")
		coinDatum.SourceCodes = strings.Split(fmt.Sprintf("%v", sourceCode), ",")
		coinDatum.Announcements = strings.Split(fmt.Sprintf("%v", announcement), ",")
		coinDatum.Slug = fmt.Sprintf("%v", slug)
		if coinDatum.Slug == "<nil>" {
			coinDatum.Slug = strings.ToLower(coinDatum.Name)
			coinDatum.Slug = strings.Replace(coinDatum.Slug, " ", "-", -1)
		}
		coinDatum.Logo = fmt.Sprintf("%v", logo)
		//fmt.Println("logo", coinDatum.Logo)
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
	}

	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists cryptos(id INTEGER, name VARCHAR, slug VARCHAR, symbol VARCHAR, logo TEXT, price REAL, vol24 REAL, date_added TEXT, explorer TEXT, bscscan TEXT, ethscan TEXT, xrpscan TEXT, bsccontract TEXT, ethcontract TEXT, xrpcontract TEXT, twitter TEXT, website TEXT, facebook TEXT, chat TEXT, message_board TEXT, technical TEXT, source_code TEXT, announcement TEXT, PRIMARY KEY(id));
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

func writeUrls(explorer string, twitter string, website string, facebook string, chat string, messageBoard string, technical string, sourceCode string, announcement string, id int64) {
	db, tx, updateExplorers := Prepare("UPDATE cryptos SET explorer =?, twitter = ?, website = ?, facebook = ?, chat = ?, message_board = ?, technical = ?, source_code = ?, announcement = ? WHERE id = ?;")
	Exec(tx, updateExplorers, explorer, twitter, website, facebook, chat, messageBoard, technical, sourceCode, announcement, fmt.Sprintf("%d", id))
	CloseDB(db)
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

func writeMap(slug string, logo string, id int64) {
	db, tx, stmt := Prepare("UPDATE cryptos SET slug = ?, logo = ? WHERE id = ?;")
	Exec(tx, stmt, slug, logo, fmt.Sprintf("%d", id))
	CloseDB(db)
}

func serializeStringList(input []string) string {
	output := ""
	for j := 0; j < len(input); j++ {
		if output == "" {
			output = input[j]

		} else {
			output += "," + input[j]
		}
	}
	return output

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
			//fmt.Println(coinDataMap.CoinDataMap[i].URLs.Explorer[j])
			//fmt.Println(strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "bscscan"))
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
			fmt.Println(coinDataMap.CoinDataMap[i].DateAdded)
			//fmt.Println(coinDataMap.CoinDataMap[i].Logo)

			writeMap(coinDataMap.CoinDataMap[i].Slug, coinDataMap.CoinDataMap[i].Logo, coinDataMap.CoinDataMap[i].Id)
		}
		twitter := serializeStringList(coinDataMap.CoinDataMap[i].URLs.Twitter)
		website := serializeStringList(coinDataMap.CoinDataMap[i].URLs.Website)
		facebook := serializeStringList(coinDataMap.CoinDataMap[i].URLs.Facebook)
		chat := serializeStringList(coinDataMap.CoinDataMap[i].URLs.Chat)
		messageBoard := serializeStringList(coinDataMap.CoinDataMap[i].URLs.MessageBoard)
		technical := serializeStringList(coinDataMap.CoinDataMap[i].URLs.Technical)
		sourceCode := serializeStringList(coinDataMap.CoinDataMap[i].URLs.SourceCode)
		announcement := serializeStringList(coinDataMap.CoinDataMap[i].URLs.Announcement)
		//fmt.Println(explorer)

		//args[1] = coinDataMap.CoinDataMap[i].Id
		//args.Explorer = explorer
		//args.Id = coinDataMap.CoinDataMap[i].Id
		//writeExplorer(explorer, coinDataMap.CoinDataMap[i].Id)
		writeUrls(explorer, twitter, website, facebook, chat, messageBoard, technical, sourceCode, announcement, coinDataMap.CoinDataMap[i].Id)
		//writeUrls(website, "website", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(facebook, "facebook", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(chat, "chat", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(messageBoard, "message_board", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(technical, "technical", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(sourceCode, "source_code", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(announcement, "announcement", coinDataMap.CoinDataMap[i].Id)
		//fmt.Println(fmt.Sprintf("%v", res))
	}

}
