package document

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	CreatedAt   time.Time
	ID          string
	Name        string
	Description string
}

type DocumentOption func(*Document) error

func NewDocument(options ...DocumentOption) (*Document, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	d := &Document{
		ID:        id.String(),
		CreatedAt: now,
	}
	for _, opt := range options {
		if err := opt(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func WithName(name string) DocumentOption {
	return func(d *Document) error {
		d.Name = name
		return nil
	}
}

func WithDescription(description string) DocumentOption {
	return func(d *Document) error {
		d.Description = description
		return nil
	}
}
