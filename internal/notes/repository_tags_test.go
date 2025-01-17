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

func TestRepository_CreateTag(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db2.RunMigrations(db, "file://../../migrations")
	require.NoError(t, err)

	repo := notes.NewRepository(db)

	t.Run("success", func(t *testing.T) {
		newTag, err := repo.CreateTag(context.Background(), "Tag text")
		require.NoError(t, err)
		require.Equal(t, "Tag text", newTag.Text)
	})
}

func TestMockNotesRepo_AddTagToNote(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = db2.RunMigrations(db, "file://../../migrations")
	require.NoError(t, err)

	repo := notes.NewRepository(db)

	t.Run("success", func(t *testing.T) {
		newTag, err := repo.CreateTag(context.Background(), "Tag text")
		require.NoError(t, err)
		note, err := repo.CreateNote(context.Background(), "Note text")
		require.NoError(t, err)

		err = repo.AddTagToNote(context.Background(), newTag.ID, note.ID)
		require.NoError(t, err)

		allNotes, err := repo.GetAllNotes(context.Background())
		require.NoError(t, err)

		require.Equal(t, newTag.Text, allNotes[0].Tags[0])
	})
}
