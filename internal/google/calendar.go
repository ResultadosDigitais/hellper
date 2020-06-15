package google

import (
	"net/http"

	"google.golang.org/api/calendar/v3"
	googleapi "google.golang.org/api/googleapi"
)

//CalendarService interfaces the Event struct from Google Calendar package
type CalendarService interface {
}

//CalendarEventsService interfaces the EventsService struct from Google Calendar package
type CalendarEventsService interface {
	Insert(string, *calendar.Event) *calendar.EventsInsertCall
}

type CalendarEventsInsertCall interface {
	// Context(context.Context) *calendar.EventsInsertCall
	Do(...googleapi.CallOption) (*calendar.Event, error)
}

//NewCalendarEventsService calls the initializer for EventsService, from the Google Calendar package
func NewCalendarEventsService(s *calendar.Service) CalendarEventsService {
	return calendar.NewEventsService(s)
}

//NewCalendarService calls the initializer for Service, from the Google Calendar package
func NewCalendarService(client *http.Client) (CalendarService, error) {
	return calendar.New(client)
}
