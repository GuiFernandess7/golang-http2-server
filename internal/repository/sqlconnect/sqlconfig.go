package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(dbname string) (*sql.DB, error) {
	fmt.Println("Trying to connect to database...")

	connectionString := "root:admin@tcp(127.0.0.1:3306)/" + dbname
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MariaDB")
	return db, nil
}
