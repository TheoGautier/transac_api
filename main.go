package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
	"transac_api/server"
)

func main() {
	logger := log.Default()
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, dbPort, dbName)
	db, err := connectToDb("postgres", connectionString)
	if err != nil {
		logger.Fatalf("Could not connect to db, err: %s", err.Error())
		return
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Fatalf("Could not retrieve port, err: %s", err.Error())
		return
	}

	s := server.MustMakeNewServer(port, db, logger)

	log.Fatalf("Server run error, err: %s", s.Start())
}

func connectToDb(driverName, dataSourceName string) (*sql.DB, error) {
	retries := 5
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = sql.Open(driverName, dataSourceName)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		return db, nil
	}
	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
