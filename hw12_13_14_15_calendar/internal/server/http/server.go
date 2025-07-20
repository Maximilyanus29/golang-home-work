package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/viper"
)

type Server struct {
	logger   Logger
	app      Application
	listener *http.Server
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, event storage.Event) error
	ListEvents(ctx context.Context) ([]storage.Event, error)
	ListEventsByDate(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
		listener: &http.Server{
			Handler: http.HandlerFunc(router),
			Addr:    net.JoinHostPort(viper.GetString("http.host"), viper.GetString("http.port")),
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.listener.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Warn(fmt.Sprintf("Ошибка сервера: %v\n", err))
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.listener.Shutdown(ctx)
	return nil
}

func router(w http.ResponseWriter, r *http.Request) {
	dateTime := time.Now().String()
	switch r.RequestURI {
	case "/":
		stringBuilder := strings.Builder{}
		stringBuilder.WriteString(r.RemoteAddr)
		stringBuilder.WriteString(" [" + dateTime + "]")
		stringBuilder.WriteString(" " + r.Method)
		stringBuilder.WriteString(" " + "HTTP/1.1")
		stringBuilder.WriteString(" " + r.RequestURI)
		stringBuilder.WriteString(" " + r.URL.Scheme)
		stringBuilder.WriteString(" " + "200")
		stringBuilder.WriteString(" " + strconv.FormatInt(r.ContentLength, 10))
		fmt.Print(stringBuilder.String())
	}
}
