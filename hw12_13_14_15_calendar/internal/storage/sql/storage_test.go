package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestSqlStorage(t *testing.T) {
	viper.SetConfigFile("../../../configs/config.toml")
	err := viper.ReadInConfig()
	require.NoError(t, err)

	ctx := context.Background()
	s := New()
	err = s.Connect(ctx, viper.GetString("storage.dsn-postgres-test"))
	require.NoError(t, err)

	_, err = s.conn.Exec("truncate events;")
	require.NoError(t, err)

	tt := time.Now()

	event := storage.Event{
		Title:            " fw",
		DateStart:        tt,
		DateEnd:          time.Now(),
		Description:      "fasfa",
		OwnerID:          1,
		TimeBeforeNotify: 0,
	}

	err = s.CreateEvent(ctx, event)
	require.NoError(t, err)

	events, err := s.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)

	require.EqualValues(t, event.Title, events[0].Title)

	events, err = s.ListEventsByDate(ctx, tt)
	require.NoError(t, err)
	require.Len(t, events, 1)

	require.EqualValues(t, event.Title, events[0].Title)

	event.Description = "wwwwww"
	err = s.UpdateEvent(ctx, event)
	require.NoError(t, err)

	events, err = s.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)

	evt := events[0]

	require.EqualValues(t, event.Description, events[0].Description)

	err = s.DeleteEvent(ctx, evt)
	require.NoError(t, err)

	events, err = s.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 0)

}
