package googlecalendar

import "hellper/internal/calendar"

type googleCalendar struct{}

//NewCalendar initialize the file storage service
func NewCalendar() calendar.Calendar {
	return new(googleCalendar)
}

//CreateCalendarEvent creates a event in Google Calendar
func (*googleCalendar) CreateCalendarEvent() {

}
