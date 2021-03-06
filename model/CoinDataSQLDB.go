package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strconv"
	"strings"
	"time"
	//"time"
)

func ReadCryptoSQLDB(id int64, DBSource string) CoinDatum {
	//fmt.Println(ConvertToISO8601(time.Now()),  "reading database...")
	db := OpenDB(DBSource)
	var err error

	rows, err := db.Query("select id, name, slug, symbol, logo, price, vol24, date_added, explorer, bscscan, ethscan, xrpscan, bsccontract, ethcontract, xrpcontract, twitter, website, facebook, chat, message_board, technical, source_code, announcement, tag, max_supply , circulating_supply, vol24, volchange24, percentchange24, percentchange7d, percentchange30d, percentchange60d, percentchange90d, market_cap, market_cap_dominance, fully_diluted_market_cap, last_price_updated_at, last_price_block_height, last_price_network_source, is_reviewed, is_rejected, is_interesting from cryptos where id = ? order by date_added desc;", fmt.Sprintf("%d", id))
	if err != nil {
		log.Fatal(err)
	}

	var coinDatum CoinDatum
	for rows.Next() {
		var slug interface{}
		var logo interface{}
		var tags interface{}
		var explorer interface{}
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
		var lastPriceUpdatedAt interface{}
		var lastPriceBlockHeight interface{}
		var lastPriceNetworkSource interface{}
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
			&announcement,
			&tags,
			&coinDatum.MaxSupply,
			&coinDatum.CirculatingSupply,
			&coinDatum.Properties.Dollar.Volume24,
			&coinDatum.Properties.Dollar.VolumeChange24,
			&coinDatum.Properties.Dollar.PercentChange24,
			&coinDatum.Properties.Dollar.PercentChange7d,
			&coinDatum.Properties.Dollar.PercentChange30d,
			&coinDatum.Properties.Dollar.PercentChange60d,
			&coinDatum.Properties.Dollar.PercentChange90d,
			&coinDatum.Properties.Dollar.MarketCap,
			&coinDatum.Properties.Dollar.MarketCapDominance,
			&coinDatum.Properties.Dollar.FullyDilutedMarketPrice,
			&lastPriceUpdatedAt,
			&lastPriceBlockHeight,
			&lastPriceNetworkSource,
			&coinDatum.IsReviewed,
			&coinDatum.IsRejected,
			&coinDatum.IsInteresting); err != nil {
			log.Fatal(err)
		}
		coinDatum.Tags = strings.Split(fmt.Sprintf("%v", tags), ",")
		coinDatum.Explorers = strings.Split(fmt.Sprintf("%v", explorer), ",")
		coinDatum.Twitters = strings.Split(fmt.Sprintf("%v", twitter), ",")
		coinDatum.Facebooks = strings.Split(fmt.Sprintf("%v", facebook), ",")
		coinDatum.Websites = strings.Split(fmt.Sprintf("%v", website), ",")
		coinDatum.MessageBoards = strings.Split(fmt.Sprintf("%v", messageBoard), ",")
		coinDatum.Chats = strings.Split(fmt.Sprintf("%v", chat), ",")
		coinDatum.Technicals = strings.Split(fmt.Sprintf("%v", technical), ",")
		coinDatum.SourceCodes = strings.Split(fmt.Sprintf("%v", sourceCode), ",")
		coinDatum.Announcements = strings.Split(fmt.Sprintf("%v", announcement), ",")
		coinDatum.Slug = fmt.Sprintf("%v", slug)
		coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = fmt.Sprintf("%v", lastPriceUpdatedAt)
		coinDatum.Properties.Dollar.LastBlock.Network = fmt.Sprintf("%v", lastPriceNetworkSource)
		if coinDatum.Slug == "<nil>" {
			coinDatum.Slug = strings.ToLower(coinDatum.Name)
			coinDatum.Slug = strings.Replace(coinDatum.Slug, " ", "-", -1)
		}
		lastPriceBlockHeightString := fmt.Sprintf("%v", lastPriceBlockHeight)
		var lastPriceBlockHeightInt int64
		if lastPriceBlockHeightString == "<nil>" {
			lastPriceBlockHeightInt = -1

		} else {
			lastPriceBlockHeightInt, _ = strconv.ParseInt(lastPriceBlockHeightString, 10, 64)
		}
		coinDatum.Properties.Dollar.LastBlock.Height = lastPriceBlockHeightInt
		coinDatum.Logo = fmt.Sprintf("%v", logo)
		//fmt.Println(ConvertToISO8601(time.Now()),  "logo", coinDatum.Logo)
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

	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", coinData))
	rows.Close()
	CloseDB(db)
	return coinDatum
}

func ReadCryptoByBSCContractSQLDB(contract string, DBSource string) CoinDatum {
	//fmt.Println(ConvertToISO8601(time.Now()),  "reading database...")
	db := OpenDB(DBSource)
	var err error
	rows, err := db.Query("select id, name, slug, symbol, logo, price, vol24, date_added, explorer, bscscan, ethscan, xrpscan, bsccontract, ethcontract, xrpcontract, twitter, website, facebook, chat, message_board, technical, source_code, announcement, tag, max_supply , circulating_supply, vol24, volchange24, percentchange24, percentchange7d, percentchange30d, percentchange60d, percentchange90d, market_cap, market_cap_dominance, fully_diluted_market_cap, last_price_updated_at, last_price_block_height, last_price_network_source, is_reviewed, is_rejected, is_interesting from cryptos where UPPER(bsccontract) LIKE UPPER('" + contract + "') order by date_added desc;")
	if err != nil {
		log.Fatal(err)
	}
	var coinDatum CoinDatum
	for rows.Next() {
		var slug interface{}
		var logo interface{}
		var tags interface{}
		var explorer interface{}
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
		var lastPriceUpdatedAt interface{}
		var lastPriceBlockHeight interface{}
		var lastPriceNetworkSource interface{}
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
			&announcement,
			&tags,
			&coinDatum.MaxSupply,
			&coinDatum.CirculatingSupply,
			&coinDatum.Properties.Dollar.Volume24,
			&coinDatum.Properties.Dollar.VolumeChange24,
			&coinDatum.Properties.Dollar.PercentChange24,
			&coinDatum.Properties.Dollar.PercentChange7d,
			&coinDatum.Properties.Dollar.PercentChange30d,
			&coinDatum.Properties.Dollar.PercentChange60d,
			&coinDatum.Properties.Dollar.PercentChange90d,
			&coinDatum.Properties.Dollar.MarketCap,
			&coinDatum.Properties.Dollar.MarketCapDominance,
			&coinDatum.Properties.Dollar.FullyDilutedMarketPrice,
			&lastPriceUpdatedAt,
			&lastPriceBlockHeight,
			&lastPriceNetworkSource,
			&coinDatum.IsReviewed,
			&coinDatum.IsRejected,
			&coinDatum.IsInteresting); err != nil {
			log.Fatal(err)
		}
		coinDatum.Tags = strings.Split(fmt.Sprintf("%v", tags), ",")
		coinDatum.Explorers = strings.Split(fmt.Sprintf("%v", explorer), ",")
		coinDatum.Twitters = strings.Split(fmt.Sprintf("%v", twitter), ",")
		coinDatum.Facebooks = strings.Split(fmt.Sprintf("%v", facebook), ",")
		coinDatum.Websites = strings.Split(fmt.Sprintf("%v", website), ",")
		coinDatum.MessageBoards = strings.Split(fmt.Sprintf("%v", messageBoard), ",")
		coinDatum.Chats = strings.Split(fmt.Sprintf("%v", chat), ",")
		coinDatum.Technicals = strings.Split(fmt.Sprintf("%v", technical), ",")
		coinDatum.SourceCodes = strings.Split(fmt.Sprintf("%v", sourceCode), ",")
		coinDatum.Announcements = strings.Split(fmt.Sprintf("%v", announcement), ",")
		coinDatum.Slug = fmt.Sprintf("%v", slug)
		coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = fmt.Sprintf("%v", lastPriceUpdatedAt)
		coinDatum.Properties.Dollar.LastBlock.Network = fmt.Sprintf("%v", lastPriceNetworkSource)
		if coinDatum.Slug == "<nil>" {
			coinDatum.Slug = strings.ToLower(coinDatum.Name)
			coinDatum.Slug = strings.Replace(coinDatum.Slug, " ", "-", -1)
		}
		lastPriceBlockHeightString := fmt.Sprintf("%v", lastPriceBlockHeight)
		var lastPriceBlockHeightInt int64
		if lastPriceBlockHeightString == "<nil>" {
			lastPriceBlockHeightInt = -1

		} else {
			lastPriceBlockHeightInt, _ = strconv.ParseInt(lastPriceBlockHeightString, 10, 64)
		}
		coinDatum.Properties.Dollar.LastBlock.Height = lastPriceBlockHeightInt
		coinDatum.Logo = fmt.Sprintf("%v", logo)
		//fmt.Println(ConvertToISO8601(time.Now()),  "logo", coinDatum.Logo)
		coinDatum.BscScan = fmt.Sprintf("%v", bscScan)
		coinDatum.EthScan = fmt.Sprintf("%v", ethScan)
		coinDatum.XrpScan = fmt.Sprintf("%v", xrpScan)
		coinDatum.BscContract = fmt.Sprintf("%v", bscContract)
		coinDatum.EthContract = fmt.Sprintf("%v", ethContract)
		coinDatum.XrpContract = fmt.Sprintf("%v", xrpContract)
		//fmt.Println(ConvertToISO8601(time.Now()),  "price: ", coinDatum.Properties.Dollar.Price)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", coinData))
	rows.Close()
	//CloseDB(db)
	return coinDatum
}

func ReadCryptosSQLDB(DBSource string) CoinData {
	//fmt.Println(ConvertToISO8601(time.Now()),  "reading database...")
	db := OpenDB(DBSource)
	var err error
	CreateCryptoTable(DBSource)
	rows, err := db.Query("select id, name, slug, symbol, logo, price, vol24, date_added, explorer, bscscan, ethscan, xrpscan, bsccontract, ethcontract, xrpcontract, twitter, website, facebook, chat, message_board, technical, source_code, announcement, tag, max_supply , circulating_supply, vol24, volchange24, percentchange24, percentchange7d, percentchange30d, percentchange60d, percentchange90d, market_cap, market_cap_dominance, fully_diluted_market_cap, last_price_updated_at, last_price_block_height, last_price_network_source, is_reviewed, is_rejected, is_interesting from cryptos order by date_added desc;")
	if err != nil {
		log.Fatal(err)
	}
	var coinData CoinData
	for rows.Next() {
		var coinDatum CoinDatum
		var name interface{}
		var slug interface{}
		var logo interface{}
		var explorer interface{}
		var twitter interface{}
		var website interface{}
		var facebook interface{}
		var chat interface{}
		var messageBoard interface{}
		var technical interface{}
		var sourceCode interface{}
		var announcement interface{}
		var tags interface{}
		var bscScan interface{}
		var ethScan interface{}
		var xrpScan interface{}
		var bscContract interface{}
		var ethContract interface{}
		var xrpContract interface{}
		var lastPriceUpdatedAt interface{}
		var lastPriceBlockHeight interface{}
		var lastPriceNetworkSource interface{}
		if err = rows.Scan(
			&coinDatum.Id,
			&name,
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
			&announcement,
			&tags,
			&coinDatum.MaxSupply,
			&coinDatum.CirculatingSupply,
			&coinDatum.Properties.Dollar.Volume24,
			&coinDatum.Properties.Dollar.VolumeChange24,
			&coinDatum.Properties.Dollar.PercentChange24,
			&coinDatum.Properties.Dollar.PercentChange7d,
			&coinDatum.Properties.Dollar.PercentChange30d,
			&coinDatum.Properties.Dollar.PercentChange60d,
			&coinDatum.Properties.Dollar.PercentChange90d,
			&coinDatum.Properties.Dollar.MarketCap,
			&coinDatum.Properties.Dollar.MarketCapDominance,
			&coinDatum.Properties.Dollar.FullyDilutedMarketPrice,
			&lastPriceUpdatedAt,
			&lastPriceBlockHeight,
			&lastPriceNetworkSource,
			&coinDatum.IsReviewed,
			&coinDatum.IsRejected,
			&coinDatum.IsInteresting); err != nil {
			log.Fatal(err)
		}
		coinDatum.Tags = strings.Split(fmt.Sprintf("%v", tags), ",")
		coinDatum.Explorers = strings.Split(fmt.Sprintf("%v", explorer), ",")
		coinDatum.Twitters = strings.Split(fmt.Sprintf("%v", twitter), ",")
		coinDatum.Facebooks = strings.Split(fmt.Sprintf("%v", facebook), ",")
		coinDatum.Websites = strings.Split(fmt.Sprintf("%v", website), ",")
		coinDatum.MessageBoards = strings.Split(fmt.Sprintf("%v", messageBoard), ",")
		coinDatum.Chats = strings.Split(fmt.Sprintf("%v", chat), ",")
		coinDatum.Technicals = strings.Split(fmt.Sprintf("%v", technical), ",")
		coinDatum.SourceCodes = strings.Split(fmt.Sprintf("%v", sourceCode), ",")
		coinDatum.Announcements = strings.Split(fmt.Sprintf("%v", announcement), ",")
		coinDatum.Slug = fmt.Sprintf("%v", slug)
		coinDatum.Name = fmt.Sprintf("%v", name)
		coinDatum.Properties.Dollar.LastBlock.TimeStamp.ISO8601 = fmt.Sprintf("%v", lastPriceUpdatedAt)
		coinDatum.Properties.Dollar.LastBlock.Network = fmt.Sprintf("%v", lastPriceNetworkSource)
		if coinDatum.Slug == "<nil>" {
			coinDatum.Slug = strings.ToLower(coinDatum.Name)
			coinDatum.Slug = strings.Replace(coinDatum.Slug, " ", "-", -1)
		}
		lastPriceBlockHeightString := fmt.Sprintf("%v", lastPriceBlockHeight)
		var lastPriceBlockHeightInt int64
		if lastPriceBlockHeightString == "<nil>" {
			lastPriceBlockHeightInt = -1

		} else {
			lastPriceBlockHeightInt, _ = strconv.ParseInt(lastPriceBlockHeightString, 10, 64)
		}
		coinDatum.Properties.Dollar.LastBlock.Height = lastPriceBlockHeightInt
		coinDatum.Logo = fmt.Sprintf("%v", logo)
		//fmt.Println(ConvertToISO8601(time.Now()),  "logo", coinDatum.Logo)
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

	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", coinData))
	rows.Close()
	//CloseDB(db)
	return coinData
}

func CreateCryptoTable(DBSource string) {
	db := OpenDB(DBSource)
	tx := TxBegin(db)
	var err error
	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists cryptos(
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name VARCHAR, slug VARCHAR, symbol VARCHAR, 
    logo TEXT, price REAL DEFAULT 0.0, vol24 REAL DEFAULT 0.0, 
    date_added TEXT, explorer VARCHAR, bscscan TEXT, 
    ethscan TEXT, xrpscan TEXT, bsccontract TEXT, 
    ethcontract TEXT, xrpcontract TEXT, twitter TEXT, 
    website TEXT, facebook TEXT, chat TEXT, 
    message_board TEXT, technical TEXT, source_code TEXT, 
    announcement TEXT, tag TEXT, max_supply REAL DEFAULT 0.0,
    circulating_supply REAL DEFAULT 0.0, volchange24 REAL DEFAULT 0.0, 
    percentchange24 REAL DEFAULT 0.0, percentchange7d REAL DEFAULT 0.0, 
    percentchange30d REAL DEFAULT 0.0, percentchange60d REAL DEFAULT 0.0, 
    percentchange90d REAL DEFAULT 0.0, market_cap REAL DEFAULT 0.0, market_cap_dominance REAL DEFAULT 0.0, 
    fully_diluted_market_cap REAL DEFAULT 0.0, last_price_updated_at TEXT,
    last_price_block_height INTEGER, last_price_network_source VARCHAR, 
    is_reviewed BIT, is_rejected BIT, is_interesting BIT);

	`); err != nil {
		log.Fatal(err)
	}
	TxCommit(tx)
	tx.Rollback()
	//time.Sleep(10 * time.Second)
	rows, err := db.Query("SELECT seq FROM sqlite_sequence WHERE name = 'cryptos'")
	rows.Next()
	var seq int64
	rows.Scan(&seq)
	rows.Close()
	//fmt.Println(ConvertToISO8601(time.Now()),  seq)
	if seq < 1000000 {
		//time.Sleep(10 * time.Second)
		if _, err = db.Exec(`
-- drop table if exists cryptos;
UPDATE sqlite_sequence
SET seq = 1000000
WHERE name = 'cryptos';`); err != nil {
			log.Fatal(err)
		}
	}
	if DBSource == "ram" {
		RamMutex.Lock()
	}
	//TxCommit(tx)
	if DBSource == "ram" {
		RamMutex.Unlock()
	}
}
func writeCrypto(id int64, name string, symbol string, dateAdded string, properties Property, tags []string, maxSupply float64, circulatingSupply float64, isReviewed bool, isRejected bool, isInteresting bool, db *sql.DB) {

	stmt := Prepare("INSERT INTO cryptos (id, name, symbol, date_added, tag, max_supply, circulating_supply, price, vol24, volchange24, percentchange24, percentchange7d, percentchange30d, percentchange60d, percentchange90d, market_cap, market_cap_dominance, fully_diluted_market_cap, last_price_updated_at, last_price_block_height, last_price_network_source, is_reviewed, is_rejected, is_interesting) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", db)
	tagString := SerializeStringList(tags)
	ExecIgnoreDuplicate(stmt, id, name, symbol, dateAdded, tagString, maxSupply, circulatingSupply, properties.Dollar.Price, properties.Dollar.Volume24, properties.Dollar.VolumeChange24, properties.Dollar.PercentChange24, properties.Dollar.PercentChange7d, properties.Dollar.PercentChange30d, properties.Dollar.PercentChange60d, properties.Dollar.PercentChange90d, properties.Dollar.MarketCap, properties.Dollar.MarketCapDominance, properties.Dollar.FullyDilutedMarketPrice, properties.Dollar.LastBlock.TimeStamp.ISO8601, properties.Dollar.LastBlock.Height, properties.Dollar.LastBlock.Network, isReviewed, isRejected, isInteresting)
	stmt.Close()
	stmt = Prepare("UPDATE cryptos SET name = ?, symbol = ?, date_added = ?, tag = ?, max_supply = ?, circulating_supply = ?, price = ?, vol24 = ?, volchange24 = ?, percentchange24 = ?, percentchange7d = ?, percentchange30d = ?, percentchange60d = ?, percentchange90d = ?, market_cap = ?, market_cap_dominance = ?, fully_diluted_market_cap = ?, last_price_updated_at = ?, last_price_block_height = ?, last_price_network_source = ?, is_reviewed = ?, is_rejected = ?, is_interesting = ? WHERE id = ? AND last_price_updated_at < ?;", db)
	Exec(stmt, name, symbol, dateAdded, tagString, maxSupply, circulatingSupply, properties.Dollar.Price, properties.Dollar.Volume24, properties.Dollar.VolumeChange24, properties.Dollar.PercentChange24, properties.Dollar.PercentChange7d, properties.Dollar.PercentChange30d, properties.Dollar.PercentChange60d, properties.Dollar.PercentChange90d, properties.Dollar.MarketCap, properties.Dollar.MarketCapDominance, properties.Dollar.FullyDilutedMarketPrice, properties.Dollar.LastBlock.TimeStamp.ISO8601, properties.Dollar.LastBlock.Height, properties.Dollar.LastBlock.Network, isReviewed, isRejected, isInteresting, id, properties.Dollar.LastBlock.TimeStamp.ISO8601)
	stmt.Close()
}
func writeCryptoFull(id int64, name string, symbol string, dateAdded string, properties Property, tags []string, maxSupply float64, circulatingSupply float64, isReviewed bool, isRejected bool, isInteresting bool, db *sql.DB) {

	stmt := Prepare("INSERT INTO cryptos (id, name, symbol, date_added, tag, max_supply, circulating_supply, price, vol24, volchange24, percentchange24, percentchange7d, percentchange30d, percentchange60d, percentchange90d, market_cap, market_cap_dominance, fully_diluted_market_cap, last_price_updated_at, last_price_block_height, last_price_network_source, is_reviewed, is_rejected, is_interesting) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", db)
	tagString := SerializeStringList(tags)
	ExecIgnoreDuplicate(stmt, id, name, symbol, dateAdded, tagString, maxSupply, circulatingSupply, properties.Dollar.Price, properties.Dollar.Volume24, properties.Dollar.VolumeChange24, properties.Dollar.PercentChange24, properties.Dollar.PercentChange7d, properties.Dollar.PercentChange30d, properties.Dollar.PercentChange60d, properties.Dollar.PercentChange90d, properties.Dollar.MarketCap, properties.Dollar.MarketCapDominance, properties.Dollar.FullyDilutedMarketPrice, properties.Dollar.LastBlock.TimeStamp.ISO8601, properties.Dollar.LastBlock.Height, properties.Dollar.LastBlock.Network, isReviewed, isRejected, isInteresting)
	stmt.Close()
	stmt = Prepare("UPDATE cryptos SET name = ?, symbol = ?, date_added = ?, tag = ?, max_supply = ?, circulating_supply = ?, price = ?, vol24 = ?, volchange24 = ?, percentchange24 = ?, percentchange7d = ?, percentchange30d = ?, percentchange60d = ?, percentchange90d = ?, market_cap = ?, market_cap_dominance = ?, fully_diluted_market_cap = ?, last_price_updated_at = ?, last_price_block_height = ?, last_price_network_source = ?, is_reviewed = ?, is_rejected = ?, is_interesting = ? WHERE id = ?;", db)
	Exec(stmt, name, symbol, dateAdded, tagString, maxSupply, circulatingSupply, properties.Dollar.Price, properties.Dollar.Volume24, properties.Dollar.VolumeChange24, properties.Dollar.PercentChange24, properties.Dollar.PercentChange7d, properties.Dollar.PercentChange30d, properties.Dollar.PercentChange60d, properties.Dollar.PercentChange90d, properties.Dollar.MarketCap, properties.Dollar.MarketCapDominance, properties.Dollar.FullyDilutedMarketPrice, properties.Dollar.LastBlock.TimeStamp.ISO8601, properties.Dollar.LastBlock.Height, properties.Dollar.LastBlock.Network, isReviewed, isRejected, isInteresting, id)
	stmt.Close()
}
func writeCryptoPrice(id int64, properties Property, db *sql.DB) {

	//stmt := Prepare("INSERT INTO cryptos (id, price) VALUES(?, ?);", db)

	//ExecIgnoreDuplicate(stmt, id, properties.Dollar.Price)
	stmt := Prepare("UPDATE cryptos SET price = ?, last_price_updated_at = ?, last_price_block_height = ?, last_price_network_source = ? WHERE id = ?;", db)
	//fmt.Println(properties.Dollar.LastBlock.Height)
	if properties.Dollar.LastBlock.Height == 0 {
		log.Fatal(id, properties)
	}
	Exec(stmt, properties.Dollar.Price, properties.Dollar.LastBlock.TimeStamp.ISO8601, properties.Dollar.LastBlock.Height, properties.Dollar.LastBlock.Network, id)
	stmt.Close()
}

func writeCryptoByBSCContract(price float64, contract string, db *sql.DB) {
	//fmt.Println(ConvertToISO8601(time.Now()),  contract)
	stmt := Prepare("UPDATE cryptos SET price = ? WHERE UPPER(bsccontract) LIKE UPPER('"+contract+"');", db)
	Exec(stmt, price)
	stmt.Close()
}

func WriteCryptosByBSCContract(data CoinData, DBSource string) {
	db := OpenDB(DBSource)
	tx := TxBegin(db)

	//fmt.Println(ConvertToISO8601(time.Now()),  "contract before db write", fmt.Sprintf("%v", data))
	for i := 0; i < len(data.CoinData); i++ {
		if i%250 == 0 {
			fmt.Println(ConvertToISO8601(time.Now()), fmt.Sprintf("contract written to db %d / %d", i, len(data.CoinData)))
		}
		if i%25 == 0 {
			time.Sleep(120 * time.Millisecond)
		}

		writeCryptoByBSCContract(data.CoinData[i].Properties.Dollar.Price, data.CoinData[i].BscContract, db)

	}
	if DBSource == "ram" {
		RamMutex.Lock()
	}
	TxCommit(tx)
	if DBSource == "ram" {
		RamMutex.Unlock()
	}

}
func WriteCryptosSQLDB(coinData CoinData, DBSource string) {

	db := OpenDB(DBSource)
	tx := TxBegin(db)
	CreateCryptoTable(DBSource)
	//CreateCryptoTable()
	fmt.Println(ConvertToISO8601(time.Now()), "writing cryptos in database...")
	for i := 0; i < len(coinData.CoinData); i++ {
		//fmt.Println(ConvertToISO8601(time.Now()),  "INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
		//	strings.Replace(coinData.CoinData[i].Name, "'", "''",-1) +"', '" +
		//	strings.Replace(coinData.CoinData[i].Symbol, "'", "''",-1) +"', '" +
		//	fmt.Sprintf("%.7+f", coinData.CoinData[i].Properties.Dollar.Price) +"', '" +
		//	fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) +"', '" +
		//	coinData.CoinData[i].DateAdded +"');")
		//
		writeCrypto(coinData.CoinData[i].Id,
			coinData.CoinData[i].Name,
			coinData.CoinData[i].Symbol,
			coinData.CoinData[i].DateAdded,
			coinData.CoinData[i].Properties,
			coinData.CoinData[i].Tags,
			coinData.CoinData[i].MaxSupply,
			coinData.CoinData[i].CirculatingSupply,
			coinData.CoinData[i].IsReviewed,
			coinData.CoinData[i].IsRejected,
			coinData.CoinData[i].IsInteresting,
			db,
		)
	}
	if DBSource == "ram" {
		RamMutex.Lock()
	}
	TxCommit(tx)
	tx.Rollback()
	if DBSource == "ram" {
		RamMutex.Unlock()
	}

}
func WriteCryptosPriceSQLDB(coinData CoinData, DBSource string) {

	CreateCryptoTable(DBSource)
	db := OpenDB(DBSource)
	tx := TxBegin(db)
	//CreateCryptoTable()
	fmt.Println(ConvertToISO8601(time.Now()), "writing cryptos in database...")
	for i := 0; i < len(coinData.CoinData); i++ {
		//fmt.Println(ConvertToISO8601(time.Now()),  "INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
		//	strings.Replace(coinData.CoinData[i].Name, "'", "''",-1) +"', '" +
		//	strings.Replace(coinData.CoinData[i].Symbol, "'", "''",-1) +"', '" +
		//	fmt.Sprintf("%.7+f", coinData.CoinData[i].Properties.Dollar.Price) +"', '" +
		//	fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) +"', '" +
		//	coinData.CoinData[i].DateAdded +"');")
		//
		writeCryptoPrice(coinData.CoinData[i].Id,
			coinData.CoinData[i].Properties,
			db,
		)
	}
	if DBSource == "ram" {
		RamMutex.Lock()
	}
	TxCommit(tx)
	tx.Rollback()
	if DBSource == "ram" {
		RamMutex.Unlock()
	}

}

func writeUrls(explorer string, twitter string, website string, facebook string, chat string, messageBoard string, technical string, sourceCode string, announcement string, id int64, db *sql.DB) {
	stmt := Prepare("UPDATE cryptos SET explorer =?, twitter = ?, website = ?, facebook = ?, chat = ?, message_board = ?, technical = ?, source_code = ?, announcement = ? WHERE id = ?;", db)
	Exec(stmt, explorer, twitter, website, facebook, chat, messageBoard, technical, sourceCode, announcement, fmt.Sprintf("%d", id))
	stmt.Close()
}

func writeExplorer(explorer string, id int64, db *sql.DB) {
	stmt := Prepare("UPDATE cryptos SET explorer = ? WHERE id = ?;", db)
	Exec(stmt, explorer, fmt.Sprintf("%d", id))
	stmt.Close()

}

func writeBscScan(bscScan string, bscContract string, id int64, db *sql.DB) {
	stmt := Prepare("UPDATE cryptos SET bscscan = ?, bsccontract = ? WHERE id = ?;", db)
	Exec(stmt, bscScan, bscContract, fmt.Sprintf("%d", id))
	stmt.Close()

}

func writeEthScan(ethScan string, ethContract string, id int64, db *sql.DB) {
	stmt := Prepare("UPDATE cryptos SET ethscan = ?, ethcontract = ? WHERE id = ?;", db)
	Exec(stmt, ethScan, ethContract, fmt.Sprintf("%d", id))
	stmt.Close()
}

func writeXrpScan(xrpScan string, xrpContract string, id int64, db *sql.DB) {
	stmt := Prepare("UPDATE cryptos SET xrpscan = ?, xrpcontract = ? WHERE id = ?;", db)
	Exec(stmt, xrpScan, xrpContract, fmt.Sprintf("%d", id))
	stmt.Close()
}

func writeMap(slug string, logo string, id int64, db *sql.DB) {
	stmt := Prepare("UPDATE cryptos SET slug = ?, logo = ? WHERE id = ?;", db)
	Exec(stmt, slug, logo, fmt.Sprintf("%d", id))
	stmt.Close()
}

func WriteCryptosMapSQLDB(coinDataMap CoinDataMap, DBSource string) {
	db := OpenDB(DBSource)
	tx := TxBegin(db)

	CreateCryptoTable(DBSource)

	for i := 0; i < len(coinDataMap.CoinDataMap); i++ {
		//fmt.Println(ConvertToISO8601(time.Now()),  "INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
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
			//fmt.Println(ConvertToISO8601(time.Now()),  coinDataMap.CoinDataMap[i].URLs.Explorer[j])
			//fmt.Println(ConvertToISO8601(time.Now()),  strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "bscscan"))
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
					writeBscScan(bscScan, bscContract, coinDataMap.CoinDataMap[i].Id, db)
				}

			}
			if strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "etherscan") {
				ethScan = coinDataMap.CoinDataMap[i].URLs.Explorer[j]
				ethContract = strings.TrimPrefix(ethScan, "https://etherscan.io/token/")
				ethContract = strings.TrimPrefix(ethContract, "https://www.etherscan.io/token/")
				writeEthScan(ethScan, ethContract, coinDataMap.CoinDataMap[i].Id, db)
			}
			if strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "xrpscan") {
				xrpScan = coinDataMap.CoinDataMap[i].URLs.Explorer[j]
				xrpContract = strings.TrimPrefix(xrpScan, "https://xrpscan.com/account/")
				xrpContract = strings.TrimPrefix(xrpContract, "https://www.xrpscan.com/account/")
				writeXrpScan(xrpScan, xrpContract, coinDataMap.CoinDataMap[i].Id, db)
			}
			//fmt.Println(ConvertToISO8601(time.Now()),  coinDataMap.CoinDataMap[i].DateAdded)
			//fmt.Println(ConvertToISO8601(time.Now()),  coinDataMap.CoinDataMap[i].Logo)

			writeMap(coinDataMap.CoinDataMap[i].Slug, coinDataMap.CoinDataMap[i].Logo, coinDataMap.CoinDataMap[i].Id, db)
		}
		twitter := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.Twitter)
		website := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.Website)
		facebook := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.Facebook)
		chat := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.Chat)
		messageBoard := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.MessageBoard)
		technical := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.Technical)
		sourceCode := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.SourceCode)
		announcement := SerializeStringList(coinDataMap.CoinDataMap[i].URLs.Announcement)
		//fmt.Println(ConvertToISO8601(time.Now()),  explorer)

		//args[1] = coinDataMap.CoinDataMap[i].Id
		//args.Explorer = explorer
		//args.Id = coinDataMap.CoinDataMap[i].Id
		//writeExplorer(explorer, coinDataMap.CoinDataMap[i].Id)
		writeUrls(explorer, twitter, website, facebook, chat, messageBoard, technical, sourceCode, announcement, coinDataMap.CoinDataMap[i].Id, db)
		//writeUrls(website, "website", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(facebook, "facebook", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(chat, "chat", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(messageBoard, "message_board", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(technical, "technical", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(sourceCode, "source_code", coinDataMap.CoinDataMap[i].Id)
		//writeUrls(announcement, "announcement", coinDataMap.CoinDataMap[i].Id)
		//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", res))
	}
	if DBSource == "ram" {
		RamMutex.Lock()
	}
	TxCommit(tx)
	tx.Rollback()
	if DBSource == "ram" {
		RamMutex.Unlock()
	}

}

func WriteCryptosFullSQLDB(coinData CoinData, DBSource string) {

	db := OpenDB(DBSource)
	tx := TxBegin(db)
	CreateCryptoTable(DBSource)
	//CreateCryptoTable()
	fmt.Println(ConvertToISO8601(time.Now()), "writing cryptos in database...")
	for i := 0; i < len(coinData.CoinData); i++ {
		//fmt.Println(ConvertToISO8601(time.Now()),  "INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
		//	strings.Replace(coinData.CoinData[i].Name, "'", "''",-1) +"', '" +
		//	strings.Replace(coinData.CoinData[i].Symbol, "'", "''",-1) +"', '" +
		//	fmt.Sprintf("%.7+f", coinData.CoinData[i].Properties.Dollar.Price) +"', '" +
		//	fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) +"', '" +
		//	coinData.CoinData[i].DateAdded +"');")
		//
		writeCryptoFull(coinData.CoinData[i].Id,
			coinData.CoinData[i].Name,
			coinData.CoinData[i].Symbol,
			coinData.CoinData[i].DateAdded,
			coinData.CoinData[i].Properties,
			coinData.CoinData[i].Tags,
			coinData.CoinData[i].MaxSupply,
			coinData.CoinData[i].CirculatingSupply,
			coinData.CoinData[i].IsReviewed,
			coinData.CoinData[i].IsRejected,
			coinData.CoinData[i].IsInteresting,
			db,
		)
		//fmt.Println(ConvertToISO8601(time.Now()),  "INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
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
		for j := 0; j < len(coinData.CoinData[i].Explorers); j++ {
			if explorer == "" {
				explorer = coinData.CoinData[i].Explorers[j]

			} else {
				explorer += "," + coinData.CoinData[i].Explorers[j]
			}
			//fmt.Println(ConvertToISO8601(time.Now()),  coinDataMap.CoinDataMap[i].URLs.Explorer[j])
			//fmt.Println(ConvertToISO8601(time.Now()),  strings.Contains(coinDataMap.CoinDataMap[i].URLs.Explorer[j], "bscscan"))
			if strings.Contains(coinData.CoinData[i].Explorers[j], "bscscan") {
				bscScan = coinData.CoinData[i].Explorers[j]

				bscContract = strings.TrimPrefix(bscScan, "https://www.bscscan.com/token/")
				bscContract = strings.TrimPrefix(bscContract, "https://bscscan.com/token/")
				bscContract = strings.TrimPrefix(bscContract, "https://bscscan.com/address/")
				bscContract = strings.TrimPrefix(bscContract, "https://www.bscscan.com/address/")
				if len(bscContract) > 42 {
					bscContract = bscContract[:42]
				}

				if strings.HasPrefix(bscContract, "0x") {
					writeBscScan(bscScan, bscContract, coinData.CoinData[i].Id, db)
				}

			}
			if coinData.CoinData[i].Id == 825 {
				writeBscScan("https://www.bscscan.com/address/0x55d398326f99059ff775485246999027b3197955", "0x55d398326f99059ff775485246999027b3197955", 825, db)
				bscContract = "0x63b7e5ae00cc6053358fb9b97b361372fba10a5e"
			}
			if strings.Contains(coinData.CoinData[i].Explorers[j], "etherscan") {
				ethScan = coinData.CoinData[i].Explorers[j]
				ethContract = strings.TrimPrefix(ethScan, "https://etherscan.io/token/")
				ethContract = strings.TrimPrefix(ethContract, "https://www.etherscan.io/token/")
				writeEthScan(ethScan, ethContract, coinData.CoinData[i].Id, db)
			}
			if strings.Contains(coinData.CoinData[i].Explorers[j], "xrpscan") {
				xrpScan = coinData.CoinData[i].Explorers[j]
				xrpContract = strings.TrimPrefix(xrpScan, "https://xrpscan.com/account/")
				xrpContract = strings.TrimPrefix(xrpContract, "https://www.xrpscan.com/account/")
				writeXrpScan(xrpScan, xrpContract, coinData.CoinData[i].Id, db)
			}
			//fmt.Println(ConvertToISO8601(time.Now()),  coinData.CoinData[i].DateAdded)
			//fmt.Println(ConvertToISO8601(time.Now()),  coinDataMap.CoinDataMap[i].Logo)

			writeMap(coinData.CoinData[i].Slug, coinData.CoinData[i].Logo, coinData.CoinData[i].Id, db)
		}
		twitter := SerializeStringList(coinData.CoinData[i].Twitters)
		website := SerializeStringList(coinData.CoinData[i].Websites)
		facebook := SerializeStringList(coinData.CoinData[i].Facebooks)
		chat := SerializeStringList(coinData.CoinData[i].Chats)
		messageBoard := SerializeStringList(coinData.CoinData[i].MessageBoards)
		technical := SerializeStringList(coinData.CoinData[i].Technicals)
		sourceCode := SerializeStringList(coinData.CoinData[i].SourceCodes)
		announcement := SerializeStringList(coinData.CoinData[i].Announcements)
		//fmt.Println(ConvertToISO8601(time.Now()),  explorer)

		//args[1] = coinDataMap.CoinDataMap[i].Id
		//args.Explorer = explorer
		//args.Id = coinDataMap.CoinDataMap[i].Id
		//writeExplorer(explorer, coinDataMap.CoinDataMap[i].Id)
		writeUrls(explorer, twitter, website, facebook, chat, messageBoard, technical, sourceCode, announcement, coinData.CoinData[i].Id, db)
	}
	if DBSource == "ram" {
		RamMutex.Lock()
	}
	TxCommit(tx)
	tx.Rollback()
	if DBSource == "ram" {
		RamMutex.Unlock()
	}

}
