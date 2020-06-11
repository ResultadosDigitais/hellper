package calendar

import (
	"context"
	"hellper/internal/model"
)

type Calendar interface {
	CreateCalendarEvent(ctx context.Context, start, end, summary, commander string, emails []string) (*model.Event, error)
}
