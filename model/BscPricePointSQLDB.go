package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"time"
	//"time"
)

func ReadPricePointsSQLDB(DBSource string) BSCPricePoints {
	//fmt.Println(ConvertToISO8601(time.Now()),  "reading database...")
	db := OpenDB(DBSource)
	var err error

	rows, err := db.Query("select id, block_height, price, network, block_date from bsc_price_points order by block_date desc;")
	if err != nil {
		log.Fatal(err)
	}

	var bscPricePoint BSCPricePoint
	var bscPricePoints BSCPricePoints
	for rows.Next() {
		if err = rows.Scan(
			&bscPricePoint.CMCId,
			&bscPricePoint.Block.Height,
			&bscPricePoint.Price,
			&bscPricePoint.Block.Network,
			&bscPricePoint.Block.TimeStamp.ISO8601); err != nil {
			log.Fatal(err)
		}
		bscPricePoints = append(bscPricePoints, bscPricePoint)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", coinData))
	rows.Close()
	CloseDB(db)
	return bscPricePoints
}

func ReadPricePointsByCMCIdSQLDB(id int64, DBSource string) BSCPricePoints {
	//fmt.Println(ConvertToISO8601(time.Now()),  "reading database...")
	db := OpenDB(DBSource)
	var err error

	rows, err := db.Query("select id, block_height, price, network, block_date from bsc_price_points where id = ? order by block_date desc;", fmt.Sprintf("%d", id))
	if err != nil {
		log.Fatal(err)
	}

	var bscPricePoint BSCPricePoint
	var bscPricePoints BSCPricePoints
	for rows.Next() {
		if err = rows.Scan(
			&bscPricePoint.CMCId,
			&bscPricePoint.Block.Height,
			&bscPricePoint.Price,
			&bscPricePoint.Block.Network,
			&bscPricePoint.Block.TimeStamp.ISO8601); err != nil {
			log.Fatal(err)
		}
		bscPricePoints = append(bscPricePoints, bscPricePoint)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", coinData))
	rows.Close()
	CloseDB(db)
	return bscPricePoints
}

func CreateBSCPricePointTable(DBSource string) {
	db := OpenDB(DBSource)
	tx := TxBegin(db)
	var err error
	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists bsc_price_points(
    id INTEGER, block_height INTEGER , price REAL, 
    network VARCHAR, block_date TEXT, PRIMARY KEY (id, block_date));
	`); err != nil {
		log.Fatal(err)
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
func writeBSCPricePoint(id int64, block BlockchainBlock, price float64, db *sql.DB) {

	stmt := Prepare("INSERT INTO bsc_price_points (id, block_height, price, network, block_date) VALUES(?, ?, ?, ?, ?);", db)
	ExecIgnoreDuplicate(stmt, id, block.Height, price, block.Network, block.TimeStamp.ISO8601)
	stmt.Close()
}

func WriteBSCPricePointsSQLDB(pricePoints BSCPricePoints, DBSource string) {

	db := OpenDB(DBSource)
	tx := TxBegin(db)
	CreateCryptoTable(DBSource)
	//CreateCryptoTable()
	fmt.Println(ConvertToISO8601(time.Now()), "writing cryptos in database...")
	for i := 0; i < len(pricePoints); i++ {
		//fmt.Println(ConvertToISO8601(time.Now()),  "INSERT INTO cryptos (name, symbol, price, vol24, date_added) VALUES ('"+
		//	strings.Replace(coinData.CoinData[i].Name, "'", "''",-1) +"', '" +
		//	strings.Replace(coinData.CoinData[i].Symbol, "'", "''",-1) +"', '" +
		//	fmt.Sprintf("%.7+f", coinData.CoinData[i].Properties.Dollar.Price) +"', '" +
		//	fmt.Sprintf("%.2f", coinData.CoinData[i].Properties.Dollar.Volume24) +"', '" +
		//	coinData.CoinData[i].DateAdded +"');")
		//
		writeBSCPricePoint(pricePoints[i].CMCId,
			pricePoints[i].Block,
			pricePoints[i].Price,
			db)
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
