package google

import (
	"net/http"

	"google.golang.org/api/calendar/v3"
)

type CalendarService interface {
}

type CalendarEventsService interface {
	Insert(string, *calendar.Event) *calendar.EventsInsertCall
}

func NewCalendarEventsService(s *calendar.Service) CalendarEventsService {
	return calendar.NewEventsService(s)
}

func NewCalendarService(client *http.Client) (CalendarService, error) {
	return calendar.New(client)
}
