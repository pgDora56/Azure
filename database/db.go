package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func CreateRepository(dbPath string) Database {
	r := Database{}
	// mysqlUrl := fmt.Sprintf(
	// 	"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	mysqlUser,
	// 	mysqlPW,
	// 	mysqlAddress,
	// 	mysqlDbname,
	// )
	// mysqlDb, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("MySQL initialization error: %s", err.Error())
	// }
	sqliteDb, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("SQLite initialization error: %s", err.Error())
	}
	r.db = sqliteDb

	err = r.migration()
	if err != nil {
		log.Fatalf("mysql migration error")
	}

	return r
}

func (r Database) migration() error {
	return r.db.AutoMigrate(
		&User{},
		// &Series{},
		// &Event{},
	)
}
