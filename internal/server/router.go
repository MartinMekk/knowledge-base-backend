package server

import (
	"github.com/gin-gonic/gin"
	"knowledge-base-backend/internal/notes"
)

func NewRouter(notesRepo notes.Repository) *gin.Engine {
	routing := gin.Default()

	h := NewHandler(notesRepo)

	routing.GET("/notes", h.GetNotesHandler)
	routing.POST("/notes", h.AddNoteHandler)

	return routing
}
