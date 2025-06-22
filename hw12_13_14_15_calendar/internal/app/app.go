package app

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

func GetRootDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath + "../"
}

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	ListEvents(ctx context.Context) ([]storage.Event, error)
	ListEventsByDate(ctx context.Context, date time.Time) ([]storage.Event, error)
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, event storage.Event) error
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, event storage.Event) error {
	return a.storage.DeleteEvent(ctx, event)
}

func (a *App) ListEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListEvents(ctx)
}

func (a *App) ListEventsByDate(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsByDate(ctx, date)
}
