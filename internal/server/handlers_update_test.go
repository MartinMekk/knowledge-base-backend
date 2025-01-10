package server_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"knowledge-base-backend/internal/notes"
	"knowledge-base-backend/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_UpdateNoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(notes.MockNotesRepo)

	h := server.NewHandler(mockRepo)

	r := gin.Default()
	r.PUT("/notes/:id", h.UpdateNoteHandler)

	makeRequest := func(noteID string, body []byte) *httptest.ResponseRecorder {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/notes/%s", noteID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("success", func(t *testing.T) {
		noteID := "1"
		requestBody := []byte(`{"text":"Updated text"}`)

		updatedNote := notes.Note{
			ID:      noteID,
			Text:    "Updated text",
			Created: time.Now(),
		}

		mockRepo.On("UpdateNote", mock.Anything, noteID, "Updated text").Return(updatedNote, nil).Once()

		w := makeRequest(noteID, requestBody)
		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		require.Equal(t, noteID, resp["id"])
		require.Equal(t, "Updated text", resp["text"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		noteID := "does-not-exist"
		requestBody := []byte(`{"text":"something"}`)

		notFoundErr := errors.New("note not found")

		mockRepo.On("UpdateNote", mock.Anything, noteID, "something").Return(notes.Note{}, notFoundErr).Once()

		w := makeRequest(noteID, requestBody)
		require.Equal(t, http.StatusNotFound, w.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		noteID := "2"
		badBody := []byte(`{invalid}`)

		w := makeRequest(noteID, badBody)
		require.Equal(t, http.StatusBadRequest, w.Code)

		mockRepo.AssertNotCalled(t, "UpdateNote", mock.Anything, noteID, mock.Anything)
	})
}
