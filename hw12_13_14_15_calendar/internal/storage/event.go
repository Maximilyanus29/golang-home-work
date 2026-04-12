package storage

import (
	"errors"
	"time"
)

var ErrDateBusy = errors.New("Собитие на эту дату уже занято")

type Event struct {
	ID               string
	Title            string
	DateTimeStart    time.Time
	DateTimeEnd      time.Time
	Description      string
	UserId           string
	TimeBeforeNotify time.Time
}
