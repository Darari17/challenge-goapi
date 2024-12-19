package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "enigma_laundry"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging the database: %v\n", err)
		return nil
	}

	fmt.Printf("Succesfully Connected!")
	return db
}