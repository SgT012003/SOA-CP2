package db

import (
	"database/sql"
	"hotel-soa/helper"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB(connectionString string) *sql.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func GetDB() *sql.DB {
	if db == nil {
		connectString := helper.GetPostgresConnectionString()
		db = connectDB(connectString)
	}
	return db
}
