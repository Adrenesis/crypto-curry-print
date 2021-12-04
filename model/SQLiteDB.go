package model

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	//"time"
)

func CloseDB(db *sql.DB) {

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

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

func ExecIgnoreDuplicate(tx *sql.Tx, stmt *sql.Stmt, args ...interface{}) {
	defer tx.Rollback()
	_, err := stmt.Exec(args...)
	s := fmt.Sprintf("%v", err)
	if err != nil && !strings.HasSuffix(s, "(1555)") {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
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
