package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	host := os.Getenv("HOST")

	fmt.Println("Trying to connect Mariadb...")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_password, host, db_port, db_name)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected Mariadb!")
	return db, nil
}
