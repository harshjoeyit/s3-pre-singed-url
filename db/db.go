package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn         = "root:mysql@tcp(localhost:3306)/test"
	maxOpenConn = 10
)

// ConnectDB establishes a connection to the database
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Set database connection pooling options
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping to ensure the connection is valid
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	// fmt.Println("Connected to MySQL successfully")
	return db, nil
}
