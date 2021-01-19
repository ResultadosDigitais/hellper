package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

const (
	StatusOpen     = "open"
	StatusCancel   = "canceled"
	StatusResolved = "resolved"
	StatusClosed   = "closed"
)

// type Incident struct {
// 	Id                      int64         `db:"id,omitempty"`
// 	Title                   string        `db:"title,omitempty"`
// 	StartTimestamp          *time.Time    `db:"start_ts,omitempty"`
// 	EndTimestamp            *time.Time    `db:"end_ts,omitempty"`
// 	IdentificationTimestamp *time.Time    `db:"identification_ts,omitempty"`
// 	Responsibility          string        `db:"responsibility,omitempty"`
// 	Team                    string        `db:"team,omitempty"`
// 	Functionality           string        `db:"functionality,omitempty"`
// 	RootCause               string        `db:"root_cause,omitempty"`
// 	CustomerImpact          sql.NullInt64 `db:"customer_impact,omitempty"`
// 	StatusPageUrl           string        `db:"status_page_url,omitempty"`
// 	PostMortemUrl           string        `db:"post_mortem_url,omitempty"`
// 	Status                  string        `db:"status,omitempty"`
// 	Product                 string        `db:"product,omitempty"`
// 	SeverityLevel           int64         `db:"severity_level,omitempty"`
// 	ChannelName             string        `db:"channel_name,omitempty"`
// 	UpdatedAt               *time.Time    `db:"updated_at,omitempty"`
// 	SnoozedUntil            sql.NullTime  `db:"snoozed_until,omitempty"`
// 	DescriptionStarted      string        `db:"description_started,omitempty"`
// 	DescriptionCancelled    string        `db:"description_cancelled,omitempty"`
// 	DescriptionResolved     string        `db:"description_resolved,omitempty"`
// 	ChannelId               string        `db:"channel_id,omitempty"`
// 	IncidentAuthor          string        `db:"incident_author_id,omitempty"`
// 	CommanderId             string        `db:"commander_id,omitempty"`
// 	CommanderEmail          string        `db:"commander_email,omitempty"`
// }

type Incident struct {
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
