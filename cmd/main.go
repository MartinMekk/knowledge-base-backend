package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	db2 "knowledge-base-backend/internal/db"
	"log"
	"net/http"
	"time"
)

type Note struct {
	ID      string    `json:"id"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
}

var db *sql.DB

func main() {
	var err error

	db, err = db2.OpenDatabase("./notes.db")
	if err != nil {
		log.Fatalf("Could not open DB: %v", err)
	}
	defer db.Close()

	err = db2.RunMigrations(db, "file://migrations")
	if err != nil {
		log.Fatalf("Could not migrate DB: %v", err)
	}

	routing := gin.Default()

	routing.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Welcome to your knowledge Base!"})
	})

	routing.GET("/notes", getNotesHandeler)

	routing.POST("/notes", addNoteHandler)

	if err := routing.Run(":8080"); err != nil {
		log.Fatalf("Failed to staret server: %v", err)
	}
}

func getNotesHandeler(c *gin.Context) {
	rows, err := db.Query("SELECT id, text, created FROM notes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Text, &note.Created); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		notes = append(notes, note)
	}
	c.JSON(http.StatusOK, notes)
}

func addNoteHandler(c *gin.Context) {
	var newNote Note
	if err := c.ShouldBindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newNote.ID = uuid.New().String()

	_, err := db.Exec("INSERT INTO notes (id, text) values (?, ?)", newNote.ID, newNote.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, newNote.ID)
}
