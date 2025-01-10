package notes

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	CreateNote(ctx context.Context, text string) (Note, error)
	GetAllNotes(ctx context.Context) ([]Note, error)
	UpdateNote(ctx context.Context, id string, newText string) (Note, error)
}

type repository struct {
	db *sql.DB
}

func (r *repository) UpdateNote(ctx context.Context, id string, newText string) (Note, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
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

func (r *repository) GetAllNotes(ctx context.Context) ([]Note, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, text, created FROM notes order by created ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Text, &note.Created); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, rows.Err()
}