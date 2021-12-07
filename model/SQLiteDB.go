package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	"time"
	//"time"
)

func CloseDB(db *sql.DB) {

	//if err := db.Close(); err != nil {
	//	log.Fatal(err)
	//}

}

var (
	DB      *sql.DB
	DBhdd   *sql.DB
	DBAlt   = false
	DBReady = false
)

func OpenDB(HDDSource bool) *sql.DB {
	var err error
	if !DBReady {
		fmt.Println("CREATING DB ########################################")
		var db *sql.DB
		var dbhdd *sql.DB

		db, err = sql.Open("sqlite", "file::memory:?cache=shared")
		DB = db
		DBReady = true
		dbhdd, err = sql.Open("sqlite", "./cryptoDB.db")
		DBhdd = dbhdd
		DBReady = true
		go LaunchSaveDeamon(300 * time.Second)
		go RewriteJsonPricesFromBitQueryEvery(20 * time.Second)
	} else {
		//DB.Close()
		//time.Sleep(2 * time.Second)
		//var db *sql.DB
		//db, err = sql.Open("sqlite", ":memory:")
		//DB = db
	}

	if err != nil {
		log.Fatal(err)
		//return
	}
	//fmt.Println("HDD DB ASKED", DBAlt)
	//fmt.Println("DB IN MEMORY", fmt.Sprintf("%v", &DB))
	if HDDSource {
		return DBhdd
	}
	return DB
}

func TxBegin(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)

	}
	return tx
}

func TxCommit(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func Prepare(query string, db *sql.DB) (stmt *sql.Stmt) {

	var err error
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Fatal(err)

	}

	return stmt
}

func Exec(stmt *sql.Stmt, args ...interface{}) {
	var err error
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)

	}
}

func ExecIgnoreDuplicate(stmt *sql.Stmt, args ...interface{}) {
	_, err := stmt.Exec(args...)
	s := fmt.Sprintf("%v", err)
	if err != nil && !strings.HasSuffix(s, "(1555)") {
		log.Fatal(err)
	}

}

func SerializeStringList(input []string) string {
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

func LaunchSaveDeamon(duration time.Duration) {
	for true {
		SaveHDDSourceEvery(duration)
	}
}

func SaveHDDSourceEvery(duration time.Duration) {
	time.Sleep(duration)
	SaveHDDSource()
}

func SaveHDDSource() {
	fmt.Println("saving.... don't touch anything...")
	DBAlt = false
	cData := ReadCryptosSQLDB(false)
	bscBalances := ReadBSCBalancesSQLDB(false)
	bscContracts := ReadBSCContractsQLDB(false)
	fmt.Println(fmt.Sprintf("%v", cData))
	DBAlt = true
	//Model.CreateCryptoTable()
	//Model.WriteCryptosSQLDB(cData)
	WriteCryptosFullSQLDB(cData, true)
	//cData = Model.ReadCryptosSQLDB()
	//fmt.Println(fmt.Sprintf("%v", cData))
	WriteBSCContractsSQLDB(bscContracts, true)
	WriteBSCBalancesSQLDB(bscBalances, true)
	DBAlt = false
	fmt.Println("saved")
}
