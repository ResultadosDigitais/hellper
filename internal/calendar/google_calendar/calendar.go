package googlecalendar

import (
	"context"
	"hellper/internal/calendar"
	"hellper/internal/google"
	googleauth "hellper/internal/google_auth"
	"hellper/internal/log"

	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendar struct {
	logger          log.Logger
	calendarService google.CalendarService
	eventsService   google.CalendarEventsService
	calendarID      string
}

//NewCalendar initialize the calendar service
func NewCalendar(
	ctx context.Context,
	logger log.Logger,
	calendarToken string,
	calendarID string,
) (calendar.Calendar, error) {
	calendarTokenBytes := []byte(calendarToken)

	gClient, err := googleauth.Struct.GetGClient(ctx, logger, calendarTokenBytes, gCalendar.CalendarScope)
	if err != nil {
		logger.Error(
			ctx,
			"googlecalendar/calendar.NewCalendar GetGClient ERROR",
			log.NewValue("error", err),
		)

		return nil, err
	}

	calendarService, err := google.NewCalendarService(gClient)
	if err != nil {
		logger.Error(
			ctx,
			"googlecalendar/calendar.NewCalendar gCalendar.New ERROR",
			log.NewValue("error", err),
		)

		return nil, err
	}

	s := calendarService.(*gCalendar.Service)
	eventsService := google.NewCalendarEventsService(s)

	calendar := googleCalendar{
		logger:          logger,
		calendarService: calendarService,
		eventsService:   eventsService,
		calendarID:      calendarID,
	}

	return &calendar, nil
}

func eventAttendee(email string, commander bool) *gCalendar.EventAttendee {
	return &gCalendar.EventAttendee{
		Email:     email,
		Organizer: commander,
	}
}

// eventDateTime receives a date-time value (formatted
// according to RFC3339) with time zone offset.
func eventDateTime(datetime string) *gCalendar.EventDateTime {
	return &gCalendar.EventDateTime{
		DateTime: datetime,
	}
}

func event(start, end, summary string, emails []string, commander string) *gCalendar.Event {
	var attendees []*gCalendar.EventAttendee
	for _, email := range emails {
		attendees = append(attendees, eventAttendee(email, false))
	}
	attendees = append(attendees, eventAttendee(commander, true))

	eventStart := eventDateTime(start)
	eventEnd := eventDateTime(end)

	return &gCalendar.Event{
		Attendees: attendees,
		Start:     eventStart,
		End:       eventEnd,
		Summary:   summary,
	}
}

func (gc *googleCalendar) insertEvent(ctx context.Context, logger log.Logger, event *gCalendar.Event) *gCalendar.Event {
	call := gc.eventsService.Insert(gc.calendarID, event)
	gcEvent, _ := gc.handleInsertEvent(ctx, logger, call)
	return gcEvent
}

func (gc *googleCalendar) handleInsertEvent(ctx context.Context, logger log.Logger, insertCall *gCalendar.EventsInsertCall) (*gCalendar.Event, error) {
	ctxEvent := insertCall.Context(ctx)
	gcEvent, err := ctxEvent.Do()
	if err != nil {
		logger.Error(
			ctx,
			"googlecalendar/calendar.insertEvent ERROR",
			log.NewValue("error", err),
		)
		return nil, err
	}
	return gcEvent, nil
}

//CreateCalendarEvent creates a event in Google Calendar
func (*googleCalendar) CreateCalendarEvent() error {
	return nil
}
