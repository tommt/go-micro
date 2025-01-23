package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPoert = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting  Auth Servic ")

	conn, _ := connectToDB()
	if conn == nil {
		log.Fatal("Failed to connect to db")
	}

	app := &Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPoert),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")

	for i := 0; i < 10; i++ {
		db, err := openDB(dsn)
		if err != nil {
			log.Println("Retrying to connect to db")
			time.Sleep(2 * time.Second)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("failed to connect to db after 10 attempts")
}
