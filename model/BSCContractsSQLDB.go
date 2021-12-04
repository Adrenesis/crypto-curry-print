package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	//"time"
)

func ReadBSCContractsQLDB() BSCContracts {
	fmt.Println("reading database...")
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}
	CreateBSCContractsTable()
	rows, err := db.Query("select * from bsccontracts;")
	if err != nil {
		log.Fatal(err)
	}
	var bscContracts BSCContracts
	for rows.Next() {
		var contract interface{}
		if err = rows.Scan(
			&contract); err != nil {
			log.Fatal(err)
		}
		contractString := fmt.Sprintf("%s", contract)
		bscContracts.Contracts = append(bscContracts.Contracts, contractString)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	CloseDB(db)
	return bscContracts
}

func CreateBSCContractsTable() {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists bsccontracts (contract VARCHAR, PRIMARY KEY(contract));
	`); err != nil {
		log.Fatal(err)
	}
}
func writeBSCContract(contract string) {
	db, tx, stmt := Prepare("INSERT INTO bsccontracts (contract) VALUES(?);")
	ExecIgnoreDuplicate(tx, stmt, contract)
	CloseDB(db)
}
func WriteBSCContractsSQLDB(bscContracts BSCContracts) {
	db, err := sql.Open("sqlite", "./cryptoDB.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	CreateBSCContractsTable()

	for i := 0; i < len(bscContracts.Contracts); i++ {
		writeBSCContract(bscContracts.Contracts[i])
	}
	CloseDB(db)

}
