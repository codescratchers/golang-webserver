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

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("database connection established")

	return db, createTable(db)
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS user (
			user_id BIGINT AUTO_INCREMENT NOT NULL UNIQUE,
			fullname VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			PRIMARY KEY (user_id),
			CONSTRAINT CK_fullname CHECK ( LENGTH(fullname) > 0 ),
			CONSTRAINT CK_email CHECK ( LENGTH(email) > 0 )
		);
	`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return err
	}

	query = `
		CREATE TABLE IF NOT EXISTS role (
		    role_id BIGINT AUTO_INCREMENT NOT NULL UNIQUE,
			role ENUM('USER', 'ADMIN') DEFAULT 'USER' NOT NULL,
			user_id BIGINT NOT NULL,
			PRIMARY KEY (role_id),
			CONSTRAINT FK_user_and_role FOREIGN KEY (user_id) REFERENCES user (user_id) ON DELETE RESTRICT
		);
	`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
