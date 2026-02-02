package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Mengambil string koneksi dari environment variable
	dsn := os.Getenv("DB_URL")

	_, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}
