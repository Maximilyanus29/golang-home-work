package sqlstorage

import (
	"context"
	"log"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	conn *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	s.conn = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.conn.NamedExec(`
	INSERT INTO events (title, date_start, date_end, descr, owner_id, time_before_notify) VALUES 
	(:title, :date_start, :date_end, :descr, :owner_id, :time_before_notify)
	`, event)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.conn.ExecContext(ctx, "update events set title = $1, date_start = $2, date_end = $3, descr = $4, owner_id = $5, time_before_notify = $6", //nolint
		event.Title,
		event.DateStart,
		event.DateEnd,
		event.Description,
		event.OwnerID,
		event.TimeBeforeNotify,
	)
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, event storage.Event) error {
	_, err := s.conn.ExecContext(ctx, "delete from events where id = $1", event.ID)
	return err
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	sql := `
	select 
		id
		, COALESCE(title, '') as title
		, date_start
		, date_end
		, COALESCE(descr, '') as descr 
		, owner_id
		, time_before_notify
	from events
	`

	events := []storage.Event{}

	err := s.conn.Select(&events, sql)

	return events, err
}

func (s *Storage) ListEventsByDate(ctx context.Context, time time.Time) ([]storage.Event, error) {
	sql := `
	select 
		id
		, COALESCE(title, '') as title
		, date_start
		, date_end
		, COALESCE(descr, '') as descr 
		, owner_id
		, time_before_notify
	from events
	where date_start = $1
	`

	events := []storage.Event{}

	err := s.conn.Select(&events, sql, time)

	return events, err
}
