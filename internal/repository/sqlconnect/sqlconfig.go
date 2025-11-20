package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	fmt.Println("Trying to connect to database...")

	connectionString := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MariaDB")
	return db, nil
}
