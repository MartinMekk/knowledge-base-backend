package notes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Repository interface {
	CreateNote(ctx context.Context, text string) (Note, error)
	GetAllNotes(ctx context.Context) ([]Note, error)
	UpdateNote(ctx context.Context, id string, newText string) (Note, error)
	CreateTag(ctx context.Context, text string) (Tag, error)
	AddTagToNote(ctx context.Context, tagId string, noteId string) error
}

type repository struct {
	db *sql.DB
}

var ErrNoteNotFound = errors.New("note not found")

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AddTagToNote(ctx context.Context, tagId string, noteId string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO note_tags (tag_id, note_id) VALUES (?, ?)`,
		tagId, noteId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateTag(ctx context.Context, text string) (Tag, error) {
	tag := Tag{ID: uuid.New().String(), Text: text}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tags (id, text) VALUES (?, ?)`,
		tag.ID, tag.Text,
	)
	if err != nil {
		return Tag{}, err
	}

	return tag, nil
}

func (r *repository) CreateNote(ctx context.Context, text string) (Note, error) {
	note := Note{ID: uuid.New().String(), Text: text, Created: time.Now()}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO notes (id, text) VALUES (?, ?)`,
		note.ID, note.Text,
	)
	if err != nil {
		return Note{}, err
	}

	row := r.db.QueryRowContext(ctx, "SELECT created FROM notes WHERE id = ?", note.ID)
	if err := row.Scan(&note.Created); err != nil {
		return Note{}, err
	}

	return note, nil
}

func (r *repository) UpdateNote(ctx context.Context, id string, newText string) (Note, error) {
	note := Note{ID: id, Text: newText}

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE notes SET text = ? WHERE id = ?`,
		newText, id,
	)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Note{}, err
	}

	if rowsAffected == 0 {
		return Note{}, ErrNoteNotFound
	}

	row := r.db.QueryRowContext(ctx, "SELECT created FROM notes WHERE id = ?", note.ID)
	if err := row.Scan(&note.Created); err != nil {
		return Note{}, err
	}

	return note, nil
}

func (r *repository) GetAllNotes(ctx context.Context) ([]Note, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT n.id AS note_id, n.text AS note_text, n.created AS note_created, GROUP_CONCAT(t.text, ',') AS tags
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		GROUP BY n.id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var tags *string
		if err := rows.Scan(&note.ID, &note.Text, &note.Created, &tags); err != nil {
			return nil, err
		}
		if tags != nil {
			note.Tags = strings.Split(*tags, ",")
		}
		notes = append(notes, note)
	}

	return notes, rows.Err()
}
