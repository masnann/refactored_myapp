package config

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

// OpenConnection membuka koneksi ke database
func OpenConnection() error {
	var err error
	db, err = setupConnection()
	if err != nil {
		return err
	}

	// Cek koneksi
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// setupConnection menyiapkan koneksi ke database
func setupConnection() (*sql.DB, error) {
	var connection = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		DBUser, DBPass, DBName, DBHost, DBPort, SSLMode)
	fmt.Println("Connection Info:", DBDriver, connection)

	db, err := sql.Open(DBDriver, connection)
	if err != nil {
		return nil, errors.New("failed to create the database connection")
	}

	return db, nil
}

// CloseConnectionDB menutup koneksi ke database
func CloseConnectionDB() {
	if db != nil {
		db.Close()
	}
}

// DBConnection mengembalikan koneksi database
func DBConnection() *sql.DB {
	return db
}
