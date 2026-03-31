package memory

import "document-metadata/pkg/document"

type DocumentStore struct {
	store *Store[document.Document]
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		store: NewStore[document.Document](),
	}
}

func (s *DocumentStore) Save(item document.Document) error {
	return s.store.Save(item.ID, item)
}

func (s *DocumentStore) FindByID(id string) (document.Document, error) {
	return s.store.FindByID(id)
}

func (s *DocumentStore) Delete(id string) error {
	return s.store.Delete(id)
}

func (s *DocumentStore) FindAll() ([]document.Document, error) {
	return s.store.FindAll()
}
