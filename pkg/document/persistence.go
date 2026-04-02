package document

import (
	"context"
	"errors"
)

var (
	ErrInvalidName = errors.New("document name cannot be empty")
)

type Persistence interface {
	Save(item Document) error
	FindAll() ([]Document, error)
	FindByID(itemId string) (Document, error)
	Delete(itemId string) error
}

type Manager struct {
	pers Persistence
}

func NewManager(pers Persistence) *Manager {
	return &Manager{pers: pers}
}

func (m *Manager) Create(ctx context.Context, name, description string) (Document, error) {
	if name == "" {
		return Document{}, ErrInvalidName
	}
	doc, err := NewDocument(
		WithName(name),
		WithDescription(description),
	)
	if err != nil {
		return Document{}, err
	}
	if err := m.pers.Save(*doc); err != nil {
		return Document{}, err
	}
	return *doc, nil
}

func (m *Manager) Get(ctx context.Context, id string) (Document, error) {
	return m.pers.FindByID(id)
}

func (m *Manager) Delete(ctx context.Context, id string) error {
	return m.pers.Delete(id)
}

func (m *Manager) FindAll(ctx context.Context) ([]Document, error) {
	return m.pers.FindAll()
}
