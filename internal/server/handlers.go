package server

import (
	"github.com/gin-gonic/gin"
	"knowledge-base-backend/internal/notes"
	"net/http"
)

type Handler struct {
	notesRepo notes.Repository
}

func NewHandler(repo notes.Repository) *Handler {
	return &Handler{notesRepo: repo}
}

func (h *Handler) GetNotesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	allNotes, err := h.notesRepo.GetAllNotes(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allNotes)
}

func (h *Handler) AddNoteHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdNote, err := h.notesRepo.CreateNote(ctx, req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdNote.ID})
}

func (h *Handler) UpdateNoteHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var req struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNote, err := h.notesRepo.UpdateNote(ctx, id, req.Text)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedNote.ID, "text": updatedNote.Text})
}
