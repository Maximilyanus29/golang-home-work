package memorystorage

import (
	"sync"

	"context"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	Events []storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Events = append(s.Events, event)
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.Events {
		if e.ID == event.ID {
			s.Events[i] = event
			return nil
		}
	}
	return storage.ErrEventNotFound
}

func (s *Storage) DeleteEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.Events {
		if e.ID == event.ID {
			s.Events = append(s.Events[:i], s.Events[i+1:]...)
			return nil
		}
	}
	return storage.ErrEventNotFound
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Events, nil
}

func (s *Storage) ListEventsByDate(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []storage.Event
	for _, e := range s.Events {
		if e.Date == date {
			events = append(events, e)
		}
	}
	return events, nil
}
