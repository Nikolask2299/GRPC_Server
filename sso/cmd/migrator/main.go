package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_"github.com/golang-migrate/migrate/v4/database/sqlite3"

	_"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to the storage directory")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to the migrations directory")
	flag.StringVar(&migrationsTable, "migrations-table", "", "Path to the migrations table")
	flag.Parse()

	if err := checkPath(storagePath, migrationsPath); err != nil {
		panic("Error: " + err.Error() + " migrations not found")
	}

	m, err := migrate.New(
		"file://" + migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable), 
	)
	
	if err != nil {
		panic(err)	
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations were applied to the database")
			return
		}

		panic(err)

	}
	
	fmt.Println("Migrations applied to the database successfully")
}




func checkPath(storagePath, migrationsPath string) error {
	if storagePath == "" {
		return errors.New("migrationsPath is required")
	}

	if migrationsPath == "" {
		return errors.New("migrationsPath is required")
	}
	return nil
}