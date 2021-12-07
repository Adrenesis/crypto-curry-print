package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	//"time"
)

func ReadBSCBalanceSQLDB(address string, contract string, HDDSource bool) BSCBalance {
	fmt.Println("reading database...")

	CreateBSCBalancesTable(HDDSource)
	db := OpenDB(HDDSource)
	var err error
	rows, err := db.Query("select address, contract, balance from bscbalances where address = ? AND contract = ?;", address, contract)
	if err != nil {
		log.Fatal(err)
	}

	var bscBalance BSCBalance
	for rows.Next() {
		var address interface{}
		var contract interface{}
		var balance float64
		if err = rows.Scan(
			&address,
			&contract,
			&balance); err != nil {
			log.Fatal(err)
		}
		bscBalance.Address = fmt.Sprintf("%s", address)
		bscBalance.Contract = fmt.Sprintf("%s", contract)
		bscBalance.Amount = balance

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//CloseDB(db)
	return bscBalance
}

func ReadBSCBalancesSQLDB(HDDSource bool) BSCBalances {
	fmt.Println("reading database...")

	CreateBSCBalancesTable(HDDSource)
	db := OpenDB(HDDSource)
	var err error
	rows, err := db.Query("select address, contract, balance from bscbalances;")
	if err != nil {
		log.Fatal(err)
	}
	var bscBalances BSCBalances
	for rows.Next() {
		var bscBalance BSCBalance
		var address interface{}
		var contract interface{}
		var balance float64
		if err = rows.Scan(
			&address,
			&contract,
			&balance); err != nil {
			log.Fatal(err)
		}
		//fmt.Println(address)
		bscBalance.Address = fmt.Sprintf("%s", address)
		bscBalance.Contract = fmt.Sprintf("%s", contract)
		bscBalance.Amount = balance
		//fmt.Println(fmt.Sprintf("%v", bscBalance))

		bscBalances.Balances = append(bscBalances.Balances, bscBalance)
		//fmt.Println(fmt.Sprintf("%v", bscBalances))
	}
	//if err = rows.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(fmt.Sprintf("%v", bscBalances))
	//CloseDB(db)
	//fmt.Println(fmt.Sprintf("%v", bscBalances))
	return bscBalances
}

func CreateBSCBalancesTable(HDDSource bool) {
	db := OpenDB(HDDSource)
	var err error
	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists bscbalances (address VARCHAR, contract VARCHAR, balance REAL, PRIMARY KEY(address, contract));
	`); err != nil {
		log.Fatal(err)
	}
}
func writeBSCBalance(address string, contract string, balance float64, db *sql.DB) {

	stmt := Prepare("INSERT INTO bscbalances (address, contract, balance) VALUES(?, ?, ?);", db)
	ExecIgnoreDuplicate(stmt, address, contract, balance)

}
func WriteBSCBalancesSQLDB(bscBalances BSCBalances, HDDSource bool) {
	CreateBSCBalancesTable(HDDSource)
	db := OpenDB(HDDSource)
	tx := TxBegin(db)

	for i := 0; i < len(bscBalances.Balances); i++ {
		writeBSCBalance(

			bscBalances.Balances[i].Address,
			bscBalances.Balances[i].Contract,
			bscBalances.Balances[i].Amount,
			db)
	}
	TxCommit(tx)
	CloseDB(db)

}
