package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

type Storage struct {
	DB *sql.DB
}

// ConnectToMySQL https://go.dev/doc/tutorial/database-access
func ConnectToMySQL(user, password, addr, dbName string) (*sql.DB, error) {
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbName,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, err
	}

	log.Println("database connection established")

	return db, err
}
