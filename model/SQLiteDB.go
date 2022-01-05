package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	"sync"
	"time"
	//"time"
)

func ConvertToISO8601(input time.Time) string {
	timeString := fmt.Sprintf("%s", input)
	timeString = timeString[:23]
	timeString += "Z"
	timeString = strings.Replace(timeString, " ", "T", -1)
	return timeString

}

func CloseDB(db *sql.DB) {

	//if err := db.Close(); err != nil {
	//	log.Fatal(err)
	//}

}

var (
	DB       *sql.DB
	DBhdd    *sql.DB
	DBprice  *sql.DB
	DBReady  = false
	RamMutex sync.Mutex
)

func OpenDB(DBSource string) *sql.DB {
	var err error
	if !DBReady {
		fmt.Println(ConvertToISO8601(time.Now()), "CREATING DB ########################################")
		var db *sql.DB
		var dbhdd *sql.DB

		db, err = sql.Open("sqlite", "file::memory:?cache=shared")
		DBprice, err = sql.Open("sqlite", "file::memory:?cache=shared")
		DB = db
		DBReady = true
		dbhdd, err = sql.Open("sqlite", "./cryptoDB.db")
		DBhdd = dbhdd
		DBReady = true
		go LaunchSaveDeamon(300 * time.Second)
		go RewriteJsonPricesFromBlockchainEvery(300 * time.Second)
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
	//fmt.Println(ConvertToISO8601(time.Now()),  "HDD DB ASKED", DBAlt)
	//fmt.Println(ConvertToISO8601(time.Now()),  "DB IN MEMORY", fmt.Sprintf("%v", &DB))
	if DBSource == "hdd" {
		return DBhdd
	} else if DBSource == "ramprice" {
		return DBprice
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
		SaveDBSourceEvery(duration)
	}
}

func SaveDBSourceEvery(duration time.Duration) {
	time.Sleep(duration)
	SaveDBSource()
}

func SaveDBSource() {
	fmt.Println(ConvertToISO8601(time.Now()), "saving.... don't touch anything...")
	cData := ReadCryptosSQLDB("ram")
	bscBalances := ReadBSCBalancesSQLDB("ram")
	bscContracts := ReadBSCContractsQLDB("ram")
	bscAdresses := ReadBSCaddressesSQLDB("ram")
	bscPricePoints := ReadPricePointsSQLDB("ram")
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", cData))
	//Model.CreateCryptoTable()
	//Model.WriteCryptosSQLDB(cData)
	WriteCryptosFullSQLDB(cData, "hdd")
	//fmt.Println(cData)
	//cData = Model.ReadCryptosSQLDB()
	//fmt.Println(ConvertToISO8601(time.Now()),  fmt.Sprintf("%v", cData))
	WriteBSCContractsSQLDB(bscContracts, "hdd")
	WriteBSCaddressesSQLDB(bscAdresses, "hdd")
	WriteBSCBalancesSQLDB(bscBalances, "hdd")
	WriteBSCPricePointsSQLDB(bscPricePoints, "hdd")
	fmt.Println(ConvertToISO8601(time.Now()), "saved")
}
