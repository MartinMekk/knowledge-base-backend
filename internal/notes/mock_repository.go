package notes

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockNotesRepo struct {
	mock.Mock
}

func (m *MockNotesRepo) CreateNote(ctx context.Context, text string) (Note, error) {
	args := m.Called(ctx, text)
	return args.Get(0).(Note), args.Error(1)
}
func (m *MockNotesRepo) GetAllNotes(ctx context.Context) ([]Note, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Note), args.Error(1)
}
func (m *MockNotesRepo) UpdateNote(ctx context.Context, id string, newText string) (Note, error) {
	args := m.Called(ctx, id, newText)
	return args.Get(0).(Note), args.Error(1)
}
