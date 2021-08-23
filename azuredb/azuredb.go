package azuredb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func execute(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func DbInitialize() *sql.DB {
	db, err := sql.Open("sqlite3", "./schedule.db")
	if err != nil {
		log.Fatal("DB file open error:", err.Error())
	}
	err = execute(db, dbInitQuery())
	if err != nil {
		log.Fatal("DB initialize error:", err.Error())
	}
	return db
}
