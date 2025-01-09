package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Note struct {
	ID      string    `json:"id"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
}

func getNotesHandeler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func addNoteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}
