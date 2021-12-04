package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	//"time"
)

func ReadBSCBalanceSQLDB(address string, contract string) BSCBalance {
	fmt.Println("reading database...")
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}
	CreateBSCBalancesTable()
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

	CloseDB(db)
	return bscBalance
}

func ReadBSCBalancesSQLDB() BSCBalances {
	fmt.Println("reading database...")
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}
	CreateBSCBalancesTable()
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
		fmt.Println(address)
		bscBalance.Address = fmt.Sprintf("%s", address)
		bscBalance.Contract = fmt.Sprintf("%s", contract)
		bscBalance.Amount = balance
		//fmt.Println(fmt.Sprintf("%v", bscBalance))

		bscBalances.Balances = append(bscBalances.Balances, bscBalance)
		fmt.Println(fmt.Sprintf("%v", bscBalances))
	}
	//if err = rows.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(fmt.Sprintf("%v", bscBalances))
	CloseDB(db)
	fmt.Println(fmt.Sprintf("%v", bscBalances))
	return bscBalances
}

func CreateBSCBalancesTable() {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists bscbalances (address VARCHAR, contract VARCHAR, balance REAL, PRIMARY KEY(address, contract));
	`); err != nil {
		log.Fatal(err)
	}
}
func writeBSCBalance(address string, contract string, balance float64) {
	db, tx, stmt := Prepare("INSERT INTO bscbalances (address, contract, balance) VALUES(?, ?, ?);")
	ExecIgnoreDuplicate(tx, stmt, address, contract, balance)
	CloseDB(db)
}
func WriteBSCBalancesSQLDB(bscBalances BSCBalances) {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	CreateBSCBalancesTable()

	for i := 0; i < len(bscBalances.Balances); i++ {
		writeBSCBalance(

			bscBalances.Balances[i].Address,
			bscBalances.Balances[i].Contract,
			bscBalances.Balances[i].Amount)
	}
	CloseDB(db)

}
