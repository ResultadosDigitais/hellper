package model

import "context"

type Repository interface {
	AddPostMortemUrl(context.Context, string, string) error
	InsertIncident(context.Context, *Incident) (int64, error)
	GetIncident(context.Context, string) (Incident, error)
	UpdateIncidentDates(context.Context, *Incident) error
	CancelIncident(context.Context, string, string) error
	CloseIncident(context.Context, *Incident) error
	ListActiveIncidents(context.Context) ([]Incident, error)
	ResolveIncident(context.Context, *Incident) error
}
