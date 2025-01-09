package main

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"knowledge-base-backend/internal/db"
	"knowledge-base-backend/internal/notes"
	"knowledge-base-backend/internal/server"
	"log"
)

func main() {
	var err error

	database, err := db.OpenDatabase("./notes.db")
	if err != nil {
		log.Fatalf("Could not open DB: %v", err)
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("Could not close DB: %v", err)
		}
	}(database)

	err = db.RunMigrations(database, "file://migrations")
	if err != nil {
		log.Fatalf("Could not migrate DB: %v", err)
	}

	notesRepo := notes.NewRepository(database)

	routing := server.NewRouter(notesRepo)

	if err := routing.Run(":8080"); err != nil {
		log.Fatalf("Failed to staret server: %v", err)
	}
}
