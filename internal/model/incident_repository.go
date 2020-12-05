package model

import "context"

// IncidentRepository wraps all database operations related to incidents
type IncidentRepository interface {
	AddPostMortemURL(context.Context, string, string) error
	InsertIncident(context.Context, *Incident) (int64, error)
	UpdateIncident(context.Context, *Incident) error
	GetIncident(context.Context, string) (Incident, error)
	UpdateIncidentDates(context.Context, *Incident) error
	CancelIncident(context.Context, *Incident) error
	CloseIncident(context.Context, *Incident) error
	ListActiveIncidents(context.Context) ([]Incident, error)
	ResolveIncident(context.Context, *Incident) error
	PauseNotifyIncident(context.Context, *Incident) error
}
