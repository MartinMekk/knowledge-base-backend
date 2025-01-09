package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(db *sql.DB) *gin.Engine {
	routing := gin.Default()

	routing.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Welcome to your knowledge Base!"})
	})

	routing.GET("/notes", getNotesHandeler(db))
	routing.POST("/notes", addNoteHandler(db))

	return routing
}
