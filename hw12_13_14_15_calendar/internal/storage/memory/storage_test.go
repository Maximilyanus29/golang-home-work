package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("test create", func(t *testing.T) {
		ctx := context.Background()
		s := New()
		event := storage.Event{
			ID:               12,
			Title:            "gasg",
			DateStart:        time.Now(),
			DateEnd:          time.Now(),
			Description:      "fgasfas",
			OwnerId:          1,
			TimeBeforeNotify: 10,
		}

		err := s.CreateEvent(ctx, event)
		require.NoError(t, err)

		list, err := s.ListEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(list), 1)
	})

	t.Run("test update", func(t *testing.T) {
		ctx := context.Background()
		s := New()
		event := storage.Event{
			ID:               1,
			Title:            "asd",
			DateStart:        time.Now(),
			DateEnd:          time.Now(),
			Description:      "fgasfas",
			OwnerId:          1,
			TimeBeforeNotify: 1,
		}

		err := s.CreateEvent(ctx, event)
		require.NoError(t, err)

		list, err := s.ListEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(list), 1)

		event.OwnerId = 3

		err = s.UpdateEvent(ctx, event)
		require.NoError(t, err)

		list, err = s.ListEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(list), 1)
		require.Equal(t, list[0].OwnerId, 3)
	})

	t.Run("test delete", func(t *testing.T) {
		ctx := context.Background()
		s := New()
		event := storage.Event{
			ID:               1,
			Title:            "1",
			DateStart:        time.Now(),
			DateEnd:          time.Now(),
			Description:      "fgasfas",
			OwnerId:          1,
			TimeBeforeNotify: 1,
		}

		err := s.CreateEvent(ctx, event)
		require.NoError(t, err)
		list, err := s.ListEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(list), 1)

		event.OwnerId = 3

		err = s.DeleteEvent(ctx, event)
		require.NoError(t, err)

		list, err = s.ListEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(list), 0)
	})

	t.Run("test list", func(t *testing.T) {
		ctx := context.Background()
		s := New()
		event := storage.Event{
			ID:               1,
			Title:            "1",
			DateStart:        time.Now(),
			DateEnd:          time.Now(),
			Description:      "fgasfas",
			OwnerId:          1,
			TimeBeforeNotify: 1,
		}

		err := s.CreateEvent(ctx, event)
		require.NoError(t, err)

		err = s.CreateEvent(ctx, event)
		require.NoError(t, err)

		err = s.CreateEvent(ctx, event)
		require.NoError(t, err)

		list, err := s.ListEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(list), 3)

	})
}
