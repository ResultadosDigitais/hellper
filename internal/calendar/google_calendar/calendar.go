package googlecalendar

import (
	"hellper/internal/calendar"

	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendar struct{}

//NewCalendar initialize the file storage service
func NewCalendar() calendar.Calendar {
	return new(googleCalendar)
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
