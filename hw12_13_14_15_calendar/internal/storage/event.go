package storage

import (
	"context"
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found")
var ErrEventDateBusy = errors.New("event date busy")
var ErrEventIDBusy = errors.New("event id busy")

type Event struct {
	ID               int       `db:"id"`
	Title            string    `db:"title"`
	DateStart        time.Time `db:"date_start"`
	DateEnd          time.Time `db:"date_end"`
	Description      string    `db:"descr"`
	OwnerID          int       `db:"owner_id"`
	TimeBeforeNotify int       `db:"time_before_notify"`
}

type StorageEvent interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, event Event) error
	ListEvents(ctx context.Context) ([]Event, error)
	ListEventsByDate(ctx context.Context, date time.Time) ([]Event, error)
}
