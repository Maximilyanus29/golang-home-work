package storage

import (
	"context"
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found")
var ErrEventBusy = errors.New("event date busy")

type Event struct {
	ID    string
	Title string
	Date  time.Time
}

type StorageEvent interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, event Event) error
	ListEvents(ctx context.Context) ([]Event, error)
	ListEventsByDate(ctx context.Context, date time.Time) ([]Event, error)
}
