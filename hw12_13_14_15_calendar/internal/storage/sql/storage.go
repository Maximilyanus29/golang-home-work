package sqlstorage

import (
	"context"
	"errors"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list = append(s.list, event)
	return nil
}

func (s *Storage) EditEvent(ctx context.Context, event *storage.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.list {
		if v.ID == event.ID {
			s.list[i] = event
			return nil
		}
	}
	errors.New("event not found")
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.list {
		if v.ID == id {
			s.list = append(s.list[:i], s.list[i+1:]...)
			return nil
		}
	}
	errors.New("event not found")
	return nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*storage.Event, len(s.list))
	copy(result, s.list)
	return result, nil
}
