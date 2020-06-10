package google

import (
	"net/http"

	"google.golang.org/api/calendar/v3"
)

//CalendarService interfaces the Event struct from Google Calendar package
type CalendarService interface {
}

//CalendarEventsService interfaces the EventsService struct from Google Calendar package
type CalendarEventsService interface {
	Insert(string, *calendar.Event) *calendar.EventsInsertCall
}

//NewCalendarEventsService calls the initializer for EventsService, from the Google Calendar package
func NewCalendarEventsService(s *calendar.Service) CalendarEventsService {
	return calendar.NewEventsService(s)
}

//NewCalendarService calls the initializer for Service, from the Google Calendar package
func NewCalendarService(client *http.Client) (CalendarService, error) {
	return calendar.New(client)
}
