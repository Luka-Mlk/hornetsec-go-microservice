package memory

import (
	"fmt"
	"sync"
)

type Store[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		items: make(map[string]T),
	}
}

func (s *Store[T]) Save(id string, item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = item
	return nil
}

func (s *Store[T]) FindByID(id string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		var zero T
		return zero, fmt.Errorf("not found")
	}
	return item, nil
}

func (s *Store[T]) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return fmt.Errorf("not found")
	}
	delete(s.items, id)
	return nil
}

func (s *Store[T]) FindAll() ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	response := []T{}
	for _, t := range s.items {
		response = append(response, t)
	}
	return response, nil
}
