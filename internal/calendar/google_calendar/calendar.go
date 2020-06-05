package googlecalendar

import (
	"context"
	"hellper/internal/calendar"
	googleauth "hellper/internal/google_auth"
	"hellper/internal/log"

	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendar struct {
	logger          log.Logger
	calendarService *gCalendar.Service
}

//NewCalendar initialize the file storage service
func NewCalendar(ctx context.Context, logger log.Logger, calendarToken string) (calendar.Calendar, error) {
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

	calendarService, err := gCalendar.New(gClient)
	if err != nil {
		logger.Error(
			ctx,
			"googlecalendar/calendar.NewCalendar gCalendar.New ERROR",
			log.NewValue("error", err),
		)

		return nil, err
	}

	calendar := googleCalendar{
		logger:          logger,
		calendarService: calendarService,
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

//CreateCalendarEvent creates a event in Google Calendar
func (*googleCalendar) CreateCalendarEvent() error {
	return nil
}
