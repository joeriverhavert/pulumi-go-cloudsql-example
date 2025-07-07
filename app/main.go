package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// Temp
	db_host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	// end Temp

	connStr := "postgres://" + db_user + ":" + db_password + "@" + db_host + "/" + db_name

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s\n", err)
	}

	defer pool.Close()

	// Create Table
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
  );`

	secondPhonesTable := `
		CREATE TABLE IF NOT EXISTS phones (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				brand TEXT NOT NULL,
				type TEXT NOT NULL
	);`

	_, err = pool.Exec(ctx, createUsersTable)
	if err != nil {
		log.Fatalf("Table creation failed: %s", err)
	}

	_, err = pool.Exec(ctx, secondPhonesTable)
	if err != nil {
		log.Fatalf("Table creation failed: %s", err)
	}

	tables, err := pool.Query(ctx, `
		SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = 'public'
		`)
	if err != nil {
		log.Fatalf("Error retrieving tables: %s", err)
	}

	for tables.Next() {
		var tableName string
		err := tables.Scan(&tableName)
		if err != nil {
			log.Fatalf("Failed to scan table: %s", err)
		}
		fmt.Println(tableName)
	}
}
