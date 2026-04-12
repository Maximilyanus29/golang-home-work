package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg ...string)
	Info(msg ...string)
	Error(msg ...string)
	Warning(msg ...string)
}

type Storage interface {
	GetEvents(ctx context.Context) ([]*storage.Event, error)
	EditEvent(ctx context.Context, event *storage.Event) error
	CreateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
}

func New(logger Logger, storage Storage) (*App, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if storage == nil {
		return nil, errors.New("storage cannot be nil")
	}
	return &App{
		logger:  logger,
		storage: storage,
	}, nil
}
func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	a.logger.Info("Creating event", "id", event.ID, "title", event.Title)

	err := a.storage.CreateEvent(ctx, event)
	if err != nil {
		a.logger.Error("Failed to create event", "error", err.Error(), "id", event.ID)
		return fmt.Errorf("create event: %w", err)
	}

	a.logger.Debug("Event created successfully", "id", event.ID)
	return nil
}

func (a *App) GetEvents(ctx context.Context) ([]*storage.Event, error) {
	a.logger.Debug("Getting all events")

	events, err := a.storage.GetEvents(ctx)
	if err != nil {
		a.logger.Error("Failed to get events", "error", err.Error())
		return nil, fmt.Errorf("get events: %w", err)
	}

	a.logger.Debug("Events retrieved", "count", strconv.Itoa(len(events)))
	return events, nil
}

func (a *App) EditEvent(ctx context.Context, event *storage.Event) error {
	a.logger.Info("Editing event", "id", event.ID)

	err := a.storage.EditEvent(ctx, event)
	if err != nil {
		a.logger.Error("Failed to edit event", "error", err.Error(), "id", event.ID)
		return fmt.Errorf("edit event: %w", err)
	}

	a.logger.Debug("Event edited successfully", "id", event.ID)
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	a.logger.Info("Deleting event", "id", id)

	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		a.logger.Error("Failed to delete event", "error", err.Error(), "id", id)
		return fmt.Errorf("delete event: %w", err)
	}

	a.logger.Debug("Event deleted successfully", "id", id)
	return nil
}
