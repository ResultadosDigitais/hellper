package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// const (
// 	StatusOpen     = "open"
// 	StatusCancel   = "canceled"
// 	StatusResolved = "resolved"
// 	StatusClosed   = "closed"
// )

type IncidentGORM struct {
	gorm.Model
	Title                string
	StartedAt            *time.Time
	EndedAt              *time.Time
	IdentifiedAt         *time.Time
	Responsibility       string
	Team                 string
	Functionality        string
	RootCause            string
	CustomerImpact       sql.NullInt64
	StatusPageURL        string
	PostMortemURL        string
	Status               string
	Product              string
	SeverityLevel        int64
	ChannelName          string
	SnoozedUntil         sql.NullTime
	DescriptionStarted   string
	DescriptionCancelled string
	DescriptionResolved  string
	ChannelID            string
	IncidentAuthor       string
	CommanderID          string
	CommanderEmail       string
}
