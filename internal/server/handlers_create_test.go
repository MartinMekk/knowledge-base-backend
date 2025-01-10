package server_test

import (
	"bytes"
	"encoding/json"
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

func TestHandler_AddNoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(notes.MockNotesRepo)

	h := server.NewHandler(mockRepo)

	r := gin.Default()
	r.POST("/notes", h.AddNoteHandler)

	makeRequest := func(body []byte) *httptest.ResponseRecorder {
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/notes"), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("success", func(t *testing.T) {
		requestBody := []byte(`{"text":"Note text"}`)

		mockRepo.On("CreateNote", mock.Anything, "Note text").Return(notes.Note{ID: "1", Text: "Note text", Created: time.Now()}, nil).Once()

		w := makeRequest(requestBody)
		require.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		require.Equal(t, "1", resp["id"])
		require.Equal(t, "Note text", resp["text"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		badBody := []byte(`{invalid}`)

		w := makeRequest(badBody)
		require.Equal(t, http.StatusBadRequest, w.Code)

		mockRepo.AssertNotCalled(t, "CreateNote", mock.Anything, mock.Anything)
	})
}
