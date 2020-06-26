package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	DB *sql.DB
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadDB(executionPath string) {
	/*
		db, err := sql.Open("sqlite3", executionPath+"/backend/database/db2.sqlite")
		checkErr(err)
	*/
	db, err := sql.Open("mysql", "ctf:1234@/ctfengine")
	log.Println(err)

	if err != nil {
		log.Panic(err)
	}
	DB = db

	log.Println("Database start")
}
