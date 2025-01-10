package notes_test

import (
	"context"
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	db2 "knowledge-base-backend/internal/db"
	"knowledge-base-backend/internal/notes"
	"testing"
	"time"
)

func TestRepository_UpdateNote_Success(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db2.RunMigrations(db, "file://../../migrations")
	require.NoError(t, err)

	repo := notes.NewRepository(db)

	originalNote, err := repo.CreateNote(context.Background(), "Original text")
	require.NoError(t, err)
	require.NotEmpty(t, originalNote.ID)

	updatedText := "Updated text"
	updatedNote, err := repo.UpdateNote(context.Background(), originalNote.ID, updatedText)
	require.NoError(t, err, "Expected successful update, got error instead")

	require.Equal(t, originalNote.ID, updatedNote.ID, "ID should be unchanged")
	require.Equal(t, updatedText, updatedNote.Text, "Text should be updated")

	require.WithinDuration(t, originalNote.Created, updatedNote.Created, time.Second, "Created timestamp should be the same")

	fetchedAll, err := repo.GetAllNotes(context.Background())
	require.NoError(t, err)
	require.Len(t, fetchedAll, 1)
	require.Equal(t, updatedText, fetchedAll[0].Text, "Text in the DB should be updated")
}

func TestRepository_UpdateNote_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db2.RunMigrations(db, "file://../../migrations")
	require.NoError(t, err)

	repo := notes.NewRepository(db)

	originalNote, err := repo.CreateNote(context.Background(), "Original text")
	require.NoError(t, err)
	require.NotEmpty(t, originalNote.ID)

	updatedText := "Updated text"
	_, err = repo.UpdateNote(context.Background(), "not an id", updatedText)
	require.NotEmpty(t, err)
	require.Equal(t, err, notes.ErrNoteNotFound)
}
