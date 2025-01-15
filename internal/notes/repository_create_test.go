package notes_test

import (
	"context"
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	db2 "knowledge-base-backend/internal/db"
	"knowledge-base-backend/internal/notes"
	"testing"
)

func TestRepository_CreateNote(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db2.RunMigrations(db, "file://../../migrations")
	require.NoError(t, err)

	repo := notes.NewRepository(db)

	t.Run("success", func(t *testing.T) {
		note, err := repo.CreateNote(context.Background(), "Note text")
		require.NoError(t, err)
		require.NotEmpty(t, note.ID)

		require.Equal(t, note.Text, "Note text")

		fetchedAll, err := repo.GetAllNotes(context.Background())
		require.NoError(t, err)
		require.Len(t, fetchedAll, 1)
		require.Equal(t, "Note text", fetchedAll[0].Text, "Text in the DB should be updated")
	})
}

func TestNewRepository_CreateNoteWithTags(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db2.RunMigrations(db, "file://../../migrations")
	require.NoError(t, err)

	repo := notes.NewRepository(db)

	t.Run("success", func(t *testing.T) {
		note, err := repo.CreateNoteWithTags(context.Background(), "Note text", []string{"Tag one", "Tag two"})
		require.NoError(t, err)
		require.NotEmpty(t, note.ID)

		require.Equal(t, note.Text, "Note text")
		require.Equal(t, note.Tags, []string{"Tag one", "Tag two"})
	})
}
