package main

import (
	"database/sql"
	"fmt"
	"log"
)

func connectDB(host string, port string, user string, password string, dbname string, queryCustom string) float64 {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected database")

	var col1 float64
	rows, err := db.Query(queryCustom)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&col1)
		fmt.Println(col1)
	}
	defer rows.Close()

	return col1
}
