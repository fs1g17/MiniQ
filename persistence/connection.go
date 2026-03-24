package persistence

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetConnectionString() string {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	pg_user := os.Getenv("POSTGRES_USER")
	pg_password := os.Getenv("POSTGRES_PASSWORD")
	pg_db := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		pg_user,
		pg_password,
		pg_db,
	)

	return psqlInfo
}

func GetConnection() {
	db, err := sql.Open("postgres", GetConnectionString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully connected to database")
}
