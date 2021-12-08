package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	//"time"
)

func ReadBSCContractsQLDB(DBSource string) BSCContracts {
	fmt.Println("reading database...")

	var err error
	CreateBSCContractsTable(DBSource)
	db := OpenDB(DBSource)
	//time.Sleep(2 * time.Second)
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

	//CloseDB(db)
	return bscContracts
}

func CreateBSCContractsTable(DBSource string) {
	db := OpenDB(DBSource)
	var err error
	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists bsccontracts (contract VARCHAR, PRIMARY KEY(contract));
	`); err != nil {
		log.Fatal(err)
	}
}
func writeBSCContract(contract string, db *sql.DB) {
	stmt := Prepare("INSERT INTO bsccontracts (contract) VALUES(?);", db)
	ExecIgnoreDuplicate(stmt, contract)

}
func WriteBSCContractsSQLDB(bscContracts BSCContracts, DBSource string) {

	fmt.Println("writing bsc contracts....")

	CreateBSCContractsTable(DBSource)
	db := OpenDB(DBSource)
	for i := 0; i < len(bscContracts.Contracts); i++ {
		stmt := Prepare("INSERT INTO bsccontracts (contract) VALUES(?);", db)
		ExecIgnoreDuplicate(stmt, bscContracts.Contracts[i])
		//writeBSCContract(bscContracts.Contracts[i])
	}
	CloseDB(db)

}
