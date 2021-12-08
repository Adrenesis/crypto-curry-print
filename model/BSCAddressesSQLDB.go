package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	//"time"
)

func ReadBSCaddressesSQLDB(DBSource string) BSCaddresses {
	fmt.Println("reading database...")

	var err error
	CreateBSCaddressesTable(DBSource)
	db := OpenDB(DBSource)
	//time.Sleep(2 * time.Second)
	rows, err := db.Query("select * from bscaddresses;")
	if err != nil {
		log.Fatal(err)
	}
	var bscaddresses BSCaddresses
	for rows.Next() {
		var addresses interface{}
		if err = rows.Scan(
			&addresses); err != nil {
			log.Fatal(err)
		}
		addressestring := fmt.Sprintf("%s", addresses)
		bscaddresses.Addresses = append(bscaddresses.Addresses, addressestring)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//CloseDB(db)
	return bscaddresses
}

func CreateBSCaddressesTable(DBSource string) {
	db := OpenDB(DBSource)
	var err error
	if _, err = db.Exec(`
-- drop table if exists cryptos;
create table if not exists bscaddresses (addresses VARCHAR, PRIMARY KEY(addresses));
	`); err != nil {
		log.Fatal(err)
	}
}
func writeBSCAddresses(addresses string, db *sql.DB) {
	stmt := Prepare("INSERT INTO bscaddresses (addresses) VALUES(?);", db)
	ExecIgnoreDuplicate(stmt, addresses)

}
func WriteBSCaddressesSQLDB(bscaddresses BSCaddresses, DBSource string) {

	fmt.Println("writing bsc addresses....")

	CreateBSCaddressesTable(DBSource)
	db := OpenDB(DBSource)
	for i := 0; i < len(bscaddresses.Addresses); i++ {
		stmt := Prepare("INSERT INTO bscaddresses (addresses) VALUES(?);", db)
		ExecIgnoreDuplicate(stmt, bscaddresses.Addresses[i])
		//writeBSCAddresses(bscaddresses.addresses[i])
	}
	CloseDB(db)

}
