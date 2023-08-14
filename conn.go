package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func connect() *sql.DB {
	var db *sql.DB
	db, err := sql.Open("mysql", "admin:SxS$5445@tcp(45.12.237.55:3306)/dbmg")
	//db, err := sql.Open("mysql", "danila:1111@tcp(127.0.0.1:3306)/dbmg")
	if err != nil {
		panic(err)
	}
	return db
}
