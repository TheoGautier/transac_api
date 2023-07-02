package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"transac_api/server"
)

func main() {
	logger := log.Default()
	host, dbPort, user, password, dbName := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbPort, user, password, dbName)
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
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
