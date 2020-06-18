package model

import "time"

type Event struct {
	EventURL string
	Start    *time.Time
	End      *time.Time
	Summary  string
}
